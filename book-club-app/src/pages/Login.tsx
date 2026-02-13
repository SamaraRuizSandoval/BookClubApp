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
} from '@ionic/react';
import { useState } from 'react';
import { useHistory } from 'react-router-dom';

import api from '../api/axios';
import { AuthTokenResponse } from '../types/auth';
import { User } from '../types/user';


type LoginProps = {
  onLoginSuccess: (authToken: AuthTokenResponse, user: User) => void;
};

export function Login(props: LoginProps) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const history = useHistory();

  const handleLogin = async () => {
    try {
      const authResponse = await api.post<AuthTokenResponse>(
        '/tokens/authentication',
        {
          username: username,
          password: password,
        },
      );

      const authToken = authResponse.data;
      const { token } = authResponse.data.auth_token;
      console.log('Received token:', token);

      localStorage.setItem('authToken', token);

      const meResponse = await api.get<User>('/me', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      const userData = meResponse.data;
      console.log('User data:', userData);

      props.onLoginSuccess(authToken, userData);
      history.replace('/');
    } catch (error) {
      console.error('Error logging in:', error);
    }
  };

  return (
    <IonPage>
      <IonContent>
        <div className="flex-center">
          <IonCard>
            <IonCardHeader>
              <IonCardTitle>Login</IonCardTitle>
            </IonCardHeader>

            <IonCardContent>
              <IonItem>
                <IonInput
                  placeholder="Username"
                  type="text"
                  value={username}
                  onIonInput={(e) => setUsername(e.detail.value!)}
                />
              </IonItem>
              <IonItem>
                <IonInput
                  placeholder="Password"
                  type="password"
                  value={password}
                  onIonInput={(e) => setPassword(e.detail.value!)}
                />
              </IonItem>
              <IonButton
                color="primary"
                shape="round"
                size="default"
                expand="full"
                className="primary-button"
                onClick={handleLogin}
              >
                Login with username
              </IonButton>
            </IonCardContent>
          </IonCard>
        </div>
      </IonContent>
    </IonPage>
  );
}
