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

export function Register() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const [usernameError, setUsernameError] = useState<string | null>(null);
  const [emailIsValid, setEmailIsValid] = useState<boolean>();
  const [emailIsTouched, setEmailTouched] = useState(false);
  const [passwordsMatch, setPasswordsMatch] = useState<boolean>();
  const [passwordIsTouched, setPasswordTouched] = useState(false);
  const history = useHistory();

  const isEmailValid = (email: string) => {
    return email.match(
      /^(?=.{1,254}$)(?=.{1,64}@)[a-zA-Z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-zA-Z0-9!#$%&'*+/=?^_`{|}~-]+)*@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/,
    );
  };

  const validateEmail = (value: string) => {
    if (value === '') {
      setEmailIsValid(undefined);
      return;
    }

    setEmailIsValid(isEmailValid(value) != null);
  };

  const markEmailTouched = () => {
    setEmailTouched(true);
  };

  const validatePasswords = (pass: string, confirm: string) => {
    if (confirm === '') {
      setPasswordsMatch(undefined);
      return;
    }

    setPasswordsMatch(pass === confirm);
  };

  const markPasswordTouched = () => {
    setPasswordTouched(true);
  };

  const handleRegister = async () => {
    if (!emailIsValid || !passwordsMatch) {
      return;
    }

    try {
      setIsLoading(true);
      setUsernameError(null); // clear previous error

      const response = await api.post('/users', {
        username,
        email,
        password,
      });

      setIsLoading(false);
      history.replace('/login', {
        message: 'Account created successfully!',
      });
      console.log('User registered:', response);
    } catch (error: any) {
      setIsLoading(false);
      if (error.response?.data?.error === 'username already taken') {
        setUsernameError('Username is already taken');
      } else if (error.response?.data?.error === 'email already in use') {
        setUsernameError('Email already in use');
      } else {
        setUsernameError('An error occurred during registration');
        console.error('Error registering user:', error);
      }
    }
  };

  return (
    <IonPage>
      <IonContent>
        <div className="flex-center">
          <IonCard>
            <IonCardHeader>
              <IonCardTitle>Create your account</IonCardTitle>
            </IonCardHeader>

            <IonCardContent>
              {usernameError && (
                <div className="error-message">{usernameError}</div>
              )}
              <IonItem lines="none">
                <IonInput
                  placeholder="Username"
                  type="text"
                  errorText="Invalid email"
                  disabled={isLoading}
                  value={username}
                  onIonInput={(e) => setUsername(e.detail.value!)}
                />
              </IonItem>
              <IonItem lines="none">
                <IonInput
                  className={`${emailIsValid && 'ion-valid'} ${emailIsValid === false && 'ion-invalid'} ${emailIsTouched && 'ion-touched'}`}
                  placeholder="Email"
                  type="email"
                  errorText="Invalid email"
                  disabled={isLoading}
                  value={email}
                  onIonInput={(e) => {
                    const value = e.detail.value!;
                    setEmail(value);
                    validateEmail(value);
                  }}
                  onIonBlur={() => {
                    markEmailTouched();
                  }}
                />
              </IonItem>
              <IonItem lines="none">
                <IonInput
                  className={`${passwordsMatch === true && 'ion-valid'} ${passwordsMatch === false && 'ion-invalid'} ${passwordIsTouched && 'ion-touched'}`}
                  placeholder="Password"
                  type="password"
                  disabled={isLoading}
                  value={password}
                  errorText="Passwords don't match"
                  onIonInput={(e) => setPassword(e.detail.value!)}
                  onIonBlur={() => {
                    markPasswordTouched();
                  }}
                />
              </IonItem>
              <IonItem lines="none">
                <IonInput
                  placeholder="Confirm password"
                  type="password"
                  errorText="Passwords don't match"
                  disabled={isLoading}
                  value={confirmPassword}
                  onIonInput={(e) => {
                    const confirmPass = e.detail.value!;
                    setConfirmPassword(confirmPass);
                    validatePasswords(password, confirmPass);
                  }}
                />
              </IonItem>
              <IonButton
                disabled={!emailIsValid || !passwordsMatch || isLoading}
                color="primary"
                shape="round"
                size="default"
                expand="full"
                className="primary-button"
                onClick={handleRegister}
              >
                Sign up with username
              </IonButton>
              <IonButton
                color="black"
                fill="clear"
                size="default"
                expand="full"
                disabled={isLoading}
                className="secondary-button"
                onClick={() => history.push('/login')}
              >
                Already have an account? Login
              </IonButton>
            </IonCardContent>
          </IonCard>
        </div>
      </IonContent>
    </IonPage>
  );
}
