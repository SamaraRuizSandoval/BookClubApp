import {
  IonPage,
  IonContent,
  IonItem,
  IonInput,
  IonButton,
  IonHeader,
  IonGrid,
  IonRow,
  IonCol,
  useIonToast,
} from '@ionic/react';
import { useState } from 'react';
import { useHistory, useLocation } from 'react-router-dom';

import { StarsBackground } from '../components/StarsBackground';
import { LandingNavBar } from '../components/landing_page/NavBar';
import { InfoPanel } from '../components/register/InfoPanel';
import '../styles/auth-forms.css';

import { api } from '../api/apiClient';

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
  const location = useLocation();
  const history = useHistory();
  const [present] = useIonToast();

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
      present({
        message: 'Account created successfully!',
        duration: 1500,
        position: 'top',
        color: 'success',
      });
      history.push('login');
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
      <IonHeader>
        <LandingNavBar />
      </IonHeader>
      <IonContent>
        <section className="midnight-bg">
          <StarsBackground />
          <div className="flex-center">
            <div>
              <IonGrid>
                <IonRow className="register-container">
                  <IonCol
                    size="12"
                    sizeMd="6"
                    className="panel-left ion-padding"
                  >
                    <InfoPanel />
                  </IonCol>
                  <IonCol size="12" sizeMd="6" className="panel-right">
                    <section
                      className="inner-right-panel"
                      aria-labelledby="form-heading"
                    >
                      <div id="formContent">
                        <span
                          className="form-eyebrow"
                          aria-label="Step: Create your account"
                        >
                          ✦ &nbsp;Create your account
                        </span>
                        <h1 className="form-heading" id="form-heading">
                          Welcome to
                          <br />
                          BookClub.
                        </h1>
                        <p className="form-sub">
                          Already have an account? <a href="/login">Sign in</a>
                        </p>
                        {usernameError && (
                          <div className="error-message">{usernameError}</div>
                        )}
                        <div className="form-grid">
                          <IonItem lines="none" className="field-group full">
                            <IonInput
                              className="custom"
                              placeholder="Username"
                              type="text"
                              errorText="Invalid email"
                              disabled={isLoading}
                              value={username}
                              onIonInput={(e) => setUsername(e.detail.value!)}
                            />
                          </IonItem>
                          <IonItem lines="none" className="field-group full">
                            <IonInput
                              className={`custom ${emailIsValid && 'ion-valid'} ${emailIsValid === false && 'ion-invalid'} ${emailIsTouched && 'ion-touched'}`}
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
                          <IonItem lines="none" className="field-group full">
                            <IonInput
                              className={`custom ${passwordsMatch === true && 'ion-valid'} ${passwordsMatch === false && 'ion-invalid'} ${passwordIsTouched && 'ion-touched'}`}
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
                          <IonItem lines="none" className="field-group full">
                            <IonInput
                              className="custom"
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
                        </div>
                        <IonButton
                          disabled={
                            !emailIsValid || !passwordsMatch || isLoading
                          }
                          size="default"
                          expand="full"
                          className="btn-submit"
                          id="submitBtn"
                          onClick={handleRegister}
                        >
                          Sign up with username
                        </IonButton>
                      </div>
                    </section>
                  </IonCol>
                </IonRow>
              </IonGrid>
            </div>
          </div>
        </section>
      </IonContent>
    </IonPage>
  );
}
