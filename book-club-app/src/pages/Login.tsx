import {
  IonPage,
  IonHeader,
  IonContent,
  IonInput,
  IonButton,
  IonItem,
} from '@ionic/react';
import { useState, useEffect } from 'react';
import { useHistory } from 'react-router-dom';

import '../styles/auth_forms.css';
import api from '../api/axios';
import { StarsBackground } from '../components/StarsBackground';
import { LandingNavBar } from '../components/landing_page/NavBar';
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

  const quotes = [
    {
      text: 'A reader lives a thousand lives before he dies. The man who never reads lives only one.',
      author: 'George R.R. Martin',
    },
    { text: 'Not all those who wander are lost.', author: 'J.R.R. Tolkien' },
    {
      text: 'There is no friend as loyal as a book.',
      author: 'Ernest Hemingway',
    },
    {
      text: 'One must always be careful of books, and what is inside them, for words have the power to change us.',
      author: 'Cassandra Clare',
    },
    { text: 'So many books, so little time.', author: 'Frank Zappa' },
  ];
  const [currentQuote, setCurrentQuote] = useState(quotes[0]);

  useEffect(() => {
    const randomQuote = quotes[Math.floor(Math.random() * quotes.length)];
    setCurrentQuote(randomQuote);
  }, []);

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
        history.replace('/admin');
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
      <IonHeader>
        <LandingNavBar />
      </IonHeader>

      <IonContent>
        <section className="midnight-bg">
          <StarsBackground />
          <div className="page-wrap">
            <div
              className="login-card"
              role="region"
              aria-label="Sign in to BookClub"
            >
              <span className="form-eyebrow">✦ &nbsp;Welcome back</span>
              <h1 className="form-heading">
                Good to see
                <br />
                you <em>again.</em>
              </h1>
              <p className="form-sub">
                Don't have an account? <a href="/register">Join for free</a>
              </p>

              {errorMessage && (
                <div className="error-message">{errorMessage}</div>
              )}
              <IonItem className="field-group full">
                <IonInput
                  className="custom"
                  placeholder="Username"
                  type="text"
                  disabled={isLoading}
                  value={username}
                  onIonInput={(e) => setUsername(e.detail.value!)}
                />
              </IonItem>
              <a href="/login" className="forgot-link">
                Forgot password?
              </a>
              <IonItem className="field-group full">
                <IonInput
                  className="custom"
                  placeholder="Password"
                  type="password"
                  disabled={isLoading}
                  value={password}
                  onIonInput={(e) => setPassword(e.detail.value!)}
                />
              </IonItem>
              <IonButton
                disabled={isLoading}
                expand="full"
                className="btn-submit"
                onClick={handleLogin}
              >
                Login with username
              </IonButton>

              <div className="or-divider" aria-hidden="true">
                <span className="or-text">or continue with</span>
              </div>
              <div className="quote-strip" aria-label="Reading quote">
                <span className="quote-mark" aria-hidden="true">
                  "
                </span>
                <p className="quote-text">{currentQuote.text}</p>
                <p className="quote-author">— {currentQuote.author}</p>
              </div>
            </div>
          </div>
        </section>
      </IonContent>
    </IonPage>
  );
}
