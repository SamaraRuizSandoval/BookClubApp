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

import api from '../api/axios';

export function Register() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleRegister = async () => {
    try {
      const response = await api.post('/users', {
        username,
        email,
        password,
      });
      console.log('User registered:', response);
    } catch (error) {
      console.error('Error registering user:', error);
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
                  placeholder="Email"
                  type="email"
                  value={email}
                  onIonInput={(e) => setEmail(e.detail.value!)}
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
                onClick={handleRegister}
              >
                Sign up with username
              </IonButton>
            </IonCardContent>
          </IonCard>
        </div>
      </IonContent>
    </IonPage>
  );
}
