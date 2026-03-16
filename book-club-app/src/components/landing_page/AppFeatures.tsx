import { IonGrid, IonRow, IonCol } from '@ionic/react';
import '../../styles/app-features.css';

export function AppFeatures() {
  return (
    <>
      <section
        className="features-section"
        id="features"
        aria-labelledby="features-heading"
      >
        {/* <div className="section-inner"> */}
        <div className="features-header reveal">
          <span className="section-label">Why BookClub</span>
          <h2 className="section-heading" id="features-heading">
            Built for serious readers.
            <br />
            Welcoming to everyone.
          </h2>
          <p className="section-sub">
            Everything you need to read more intentionally and connect with
            people who share your taste.
          </p>
        </div>
        <IonGrid className='class="features-grid"'>
          <IonRow>
            <IonCol size="12" sizeMd="6">
              <div className="feature-card reveal reveal-delay-1">
                <span className="feature-num" aria-hidden="true">
                  01
                </span>
                <span className="feature-icon" aria-hidden="true">
                  📚
                </span>
                <h3 className="feature-title">Track Your Books</h3>
                <p className="feature-body">
                  Keep a personal shelf of what you're reading, finished, or
                  want to read next.
                </p>
              </div>
            </IonCol>
            <IonCol size="12" sizeMd="6">
              <div className="feature-card reveal reveal-delay-2">
                <span className="feature-num" aria-hidden="true">
                  02
                </span>
                <span className="feature-icon" aria-hidden="true">
                  💬
                </span>
                <h3 className="feature-title">Share Your Thoughts</h3>
                <p className="feature-body">
                  Write reviews and reactions that your community can see and
                  respond to.
                </p>
              </div>
            </IonCol>
          </IonRow>
          <IonRow>
            <IonCol size="12" sizeMd="6">
              <div className="feature-card reveal reveal-delay-3">
                <span className="feature-num" aria-hidden="true">
                  03
                </span>
                <span className="feature-icon" aria-hidden="true">
                  🔍
                </span>
                <h3 className="feature-title">Discover New Books</h3>
                <p className="feature-body">
                  Browse community picks and personalized suggestions based on
                  what you love.
                </p>
              </div>
            </IonCol>
            <IonCol size="12" sizeMd="6">
              <div className="feature-card  reveal reveal-delay-4">
                <span className="feature-num" aria-hidden="true">
                  04
                </span>
                <span className="feature-icon" aria-hidden="true">
                  👥
                </span>
                <h3 className="feature-title">Join the Conversation</h3>
                <p className="feature-body">
                  Engage with readers who share your taste in a warm,
                  judgment-free space.
                </p>
              </div>
            </IonCol>
          </IonRow>
        </IonGrid>
        {/* </div> */}
      </section>
    </>
  );
}
