import { IonContent } from '@ionic/react';
import '../../styles/app_features.css';

export function AppFeatures() {
  return (
    <>
      <section
        className="features-section"
        id="features"
        aria-labelledby="features-heading"
      >
        <div className="section-inner">
          <div className="features-header ">
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
        </div>
      </section>
    </>
  );
}
