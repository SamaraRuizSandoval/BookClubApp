import {
  IonMenu,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonContent,
  IonMenuToggle,
  IonSpinner,
} from '@ionic/react';
import { useLocation } from 'react-router-dom';

import { useUserStats } from '../context/UserStatsContext';

import '../styles/left-menu.css';
import '../global.css';
import { key } from 'ionicons/icons';

const navItems = [
  { to: '/app', icon: '🔍', label: 'Browse Books' },
  { to: '/app/library', icon: '🏠', label: 'My Shelf' },
];

type BookStatusKey = 'reading' | 'wishlist' | 'completed';
const bookCollections: {
  to: string;
  style: React.CSSProperties;
  label: string;
  key: BookStatusKey;
  id: string;
}[] = [
  {
    to: '/app/reading',
    style: { background: 'var(--reading)' },
    label: 'Reading',
    key: 'reading',
    id: 'count-reading',
  },
  {
    to: '/app/wishlist',
    style: { background: 'var(--want)' },
    label: 'Wishlist',
    key: 'wishlist',
    id: 'count-want',
  },
  {
    to: '/app/completed',
    style: { background: 'var(--done)' },
    label: 'Completed',
    key: 'completed',
    id: 'count-done',
  },
];

export function LeftMenu() {
  const location = useLocation();
  const { stats, refreshStats } = useUserStats();

  const renderNavItem = (to: string, icon: any, label: string) => {
    const active = location.pathname.startsWith(to);
    return (
      <IonMenuToggle key={to} autoHide={false}>
        <a className={`nav-item ${active ? 'active' : ''}`} href={to}>
          <span className="nav-icon">{icon}</span>
          {label}
        </a>
      </IonMenuToggle>
    );
  };

  return (
    <IonMenu
      className="sidebar"
      contentId="main-content"
      type="overlay"
      side="start"
    >
      <IonHeader>
        <IonToolbar>
          <IonTitle>
            <a href="/" className="topbar-logo" aria-label="BookClub home">
              <span className="logo-mark" aria-hidden="true">
                📚
              </span>
              BookClub
            </a>
          </IonTitle>
        </IonToolbar>
      </IonHeader>

      <IonContent fullscreen className="sidebar-content">
        {/* Top nav items */}
        <div className="sidebar-section">
          {navItems.map(({ to, icon, label }) =>
            renderNavItem(to, icon, label),
          )}
        </div>

        <div className="sidebar-divider" />

        {/* Book collections */}
        <div className="sidebar-section">
          <span className="sidebar-label">My Collections</span>
          {bookCollections.map(({ to, style, label, key, id }) => {
            const active = location.pathname.startsWith(to);
            return (
              <IonMenuToggle key={to} autoHide={false}>
                <a
                  className={`collection-mini ${active ? 'active' : ''}`}
                  href={to}
                >
                  <span className="coll-dot" style={style}></span>
                  {label}
                  <span className="nav-count gold" id={id}>
                    {stats ? stats?.[key] : <IonSpinner name="dots" />}
                  </span>
                </a>
              </IonMenuToggle>
            );
          })}
        </div>
      </IonContent>
    </IonMenu>
  );
}
