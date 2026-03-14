import { IonPage, IonHeader, IonContent } from '@ionic/react';
import { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';

import { AppFeatures } from '../components/landing_page/AppFeatures';
import { Hero } from '../components/landing_page/Hero';
import { LandingNavBar } from '../components/landing_page/NavBar';
import { useScrollReveal } from '../hooks/useScrollReveal';
import './LandingPage.css';

export function LandingPage() {
  useScrollReveal();
  const history = useHistory();
  const [showFooter, setShowFooter] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      const scrollY = window.scrollY;
      // Show footer after scrolling 200px (adjust as needed)
      setShowFooter(scrollY > 200);
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <IonPage>
      {/* Nav */}
      <IonHeader>
        <LandingNavBar />
      </IonHeader>

      <IonContent>
        <Hero />

        <AppFeatures />

        <div className="footer" id="about">
          <div className="footer-inner">
            <a href="/" className="footer-logo" aria-label="BookClub home">
              <span className="logo-mark" aria-hidden="true">
                📚
              </span>
              BookClub
            </a>
            <nav className="footer-links" aria-label="Footer navigation">
              <a href="/">About</a>
              <a
                href="https://github.com/SamaraRuizSandoval/BookClubApp"
                target="_blank"
                rel="noopener"
              >
                GitHub
              </a>
              <a href="/">Privacy Policy</a>
              <a href="/">Contact</a>
            </nav>
            <p className="footer-copy">
              📚 Books are better together. · © 2026 BookClub
            </p>
          </div>
        </div>
      </IonContent>
    </IonPage>
  );
}
