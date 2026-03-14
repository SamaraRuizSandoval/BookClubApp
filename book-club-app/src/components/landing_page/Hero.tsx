import { IonGrid, IonRow, IonCol, IonButton } from '@ionic/react';
import { useHistory } from 'react-router-dom';

import { BookAnimation } from './BooksAnimation';

import '../../styles/hero.css';
import { StarsBackground } from '../StarsBackground';
export function Hero() {
  const history = useHistory();
  return (
    <>
      <section className="hero">
        <IonGrid>
          <StarsBackground />
          <IonRow className="ion-align-items-center">
            <IonCol size="12" sizeMd="6" className="hero-text">
              <div className="hero-eyebrow">
                ✨ Free to join · No ads · Just books
              </div>
              <h1 className="hero-headline">
                Books are
                <br />
                <em>better together.</em>
              </h1>
              <p className="hero-sub">
                <span>
                  Track what you read, share what you love, and discover your
                  next favorite book — with a community that gets you.
                </span>
              </p>
              <IonButton
                className="btn-primary"
                onClick={() => history.push('/register')}
              >
                Join for Free
              </IonButton>
              <IonButton
                className="hero-login"
                fill="clear"
                onClick={() => history.push('/login')}
              >
                Already a member? Log in
              </IonButton>
            </IonCol>
            <IonCol size="12" sizeMd="6">
              <BookAnimation></BookAnimation>
            </IonCol>
          </IonRow>
        </IonGrid>
      </section>
    </>
  );
}
