import {
  IonMenu,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonContent,
  IonMenuToggle,
} from '@ionic/react';
import { useLocation } from 'react-router-dom';
import '../styles/left-menu.css';
import '../global.css';

const navItems = [
  { to: '/home', icon: '🔍', label: 'Browse Books' },
  { to: '/my-shelf', icon: '🏠', label: 'My Shelf' },
];

const bookCollections = [
  {
    to: '/reading',
    style: { background: 'var(--reading)' },
    label: 'Reading',
    id: 'count-reading',
  },
  {
    to: '/wishlist',
    style: { background: 'var(--want)' },
    label: 'Wishlist',
    id: 'count-want',
  },
  {
    to: '/completed',
    style: { background: 'var(--done)' },
    label: 'Completed',
    id: 'count-done',
  },
];

export function LeftMenu() {
  const location = useLocation();

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
          {bookCollections.map(({ to, style, label, id }) => {
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
                    0
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
