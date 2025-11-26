import {
  IonSplitPane,
  IonMenu,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonContent,
  IonList,
  IonItem,
  IonMenuToggle,
  IonLabel,
  IonRouterOutlet,
} from '@ionic/react';
import React from 'react';
import {
  Switch,
  Route,
  useLocation,
  useParams,
  useRouteMatch,
} from 'react-router-dom';

import { allBooks, BookStatus } from './../bookData';
import { Page } from './../pages/Page';

type Props = {
  sectionKey: BookStatus; // 'reading' | 'wishlist' | 'completed'
  sectionTitle: string; // Human-readable title
};

export function BooksSectionLayout({ sectionKey, sectionTitle }: Props) {
  const { path, url } = useRouteMatch();

  const books = allBooks.filter((b) => b.status === sectionKey);
  const contentId = `${sectionKey}-content`;
  const menuId = `${sectionKey}-books-menu`;

  function SectionBookPage() {
    const { bookId } = useParams<{ bookId: string }>();
    const book = books.find((b) => b.id === bookId);
    return <Page title={book ? book.title : sectionTitle} />;
  }

  return (
    <IonSplitPane when="lg" contentId={contentId}>
      {/* inner left menu for this section */}
      <IonMenu
        contentId={contentId}
        side="start"
        type="overlay"
        menuId={menuId}
      >
        <IonHeader>
          <IonToolbar>
            <IonTitle>{sectionTitle}</IonTitle>
          </IonToolbar>
        </IonHeader>
        <IonContent>
          <BooksMenuList baseUrl={url} books={books} />
        </IonContent>
      </IonMenu>

      {/* right side: section content */}
      <IonRouterOutlet id={contentId}>
        <Switch>
          <Route
            exact
            path={path}
            component={() => <Page title={sectionTitle} />}
          />
          <Route path={`${path}/:bookId`} component={SectionBookPage} />
        </Switch>
      </IonRouterOutlet>
    </IonSplitPane>
  );
}

// Small pure component just for the list, so we can reuse it if needed.
type MenuListProps = {
  baseUrl: string;
  books: { id: string; title: string }[];
};

function BooksMenuList({ baseUrl, books }: MenuListProps) {
  const location = useLocation();

  return (
    <IonList>
      {books.map((book) => {
        const to = `${baseUrl}/${book.id}`;
        const active = location.pathname.startsWith(to);
        return (
          <IonMenuToggle key={book.id} autoHide={false}>
            <IonItem
              button
              detail={false}
              color={active ? 'primary' : undefined}
              routerLink={to}
              routerDirection="root"
            >
              <IonLabel>{book.title}</IonLabel>
            </IonItem>
          </IonMenuToggle>
        );
      })}
    </IonList>
  );
}
