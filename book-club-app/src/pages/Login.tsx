import {
  IonPage,
  IonContent,
  IonCard,
  IonCardHeader,
  IonCardContent,
  IonCardTitle,
  IonItem,
  IonInput,
  IonButton,
  IonToast,
} from '@ionic/react';
import { useState } from 'react';
import { useHistory, useLocation } from 'react-router-dom';

import api from '../api/axios';
import { useAuth } from '../context/AuthContext';
import { AuthTokenResponse } from '../types/auth';
import { User } from '../types/user';

type LoginLocationState = {
  message?: string;
};

export function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const { login } = useAuth();
  const history = useHistory();
  const [isLoading, setIsLoading] = useState(false);

  const location = useLocation<LoginLocationState>();
  const successMessage = (location.state as LoginLocationState | undefined)
    ?.message;
  const [showToast, setShowToast] = useState(!!successMessage);

  const handleLogin = async () => {
    setErrorMessage(null); // Clear previous error message

    try {
      setIsLoading(true);

      const authResponse = await api.post<AuthTokenResponse>(
        '/tokens/authentication',
        {
          username: username,
          password: password,
        },
      );

      const authToken = authResponse.data.auth_token.token;

      localStorage.setItem('authToken', authToken);

      const userData = await api.get<User>('/me', {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });

      login(authToken, userData);
      setIsLoading(false);
      if (userData.data.role === 'admin') {
        history.replace('/dashboard');
      } else {
        history.replace('/home');
      }
    } catch (error) {
      console.error('Error logging in:', error);
      setErrorMessage('Login failed. Invalid username or password.');
      setIsLoading(false);
    }
  };

  return (
    <IonPage>
      <IonToast
        isOpen={showToast}
        message={successMessage}
        position="top"
        duration={3000}
        color="success"
        onDidDismiss={() => setShowToast(false)}
      />
      <IonContent>
        <div className="flex-center">
          <IonCard>
            <IonCardHeader>
              <IonCardTitle>Login</IonCardTitle>
            </IonCardHeader>

            <IonCardContent>
              {errorMessage && (
                <div className="error-message">{errorMessage}</div>
              )}
              <IonItem>
                <IonInput
                  placeholder="Username"
                  type="text"
                  disabled={isLoading}
                  value={username}
                  onIonInput={(e) => setUsername(e.detail.value!)}
                />
              </IonItem>
              <IonItem>
                <IonInput
                  placeholder="Password"
                  type="password"
                  disabled={isLoading}
                  value={password}
                  onIonInput={(e) => setPassword(e.detail.value!)}
                />
              </IonItem>
              <IonButton
                disabled={isLoading}
                color="primary"
                shape="round"
                size="default"
                expand="full"
                className="primary-button"
                onClick={handleLogin}
              >
                Login with username
              </IonButton>

              <IonButton
                disabled={isLoading}
                color="black"
                fill="clear"
                size="default"
                expand="full"
                className="secondary-button"
                onClick={() => history.push('/register')}
              >
                Don't have an account? Register
              </IonButton>
            </IonCardContent>
          </IonCard>
        </div>
      </IonContent>
    </IonPage>
  );
}
