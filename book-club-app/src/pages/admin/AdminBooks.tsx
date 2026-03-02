import {
  IonContent,
  IonPage,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonText,
  IonSpinner,
  IonButton,
  IonIcon,
} from '@ionic/react';
import { addOutline } from 'ionicons/icons';
import { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';

import api from '../../api/axios';
import { BookGrid } from '../../components/BooksGrid';
import { Book, BookResponse } from '../../types/book';

export function AdminBooks() {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(true);
  const history = useHistory();

  useEffect(() => {
    async function fetchBooks() {
      try {
        const response = await api.get<BookResponse>('/books');
        setBooks(response.data.books);
        console.log(response.data);
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    }
    fetchBooks();
  }, []);

  // TODO: We can add caching
  // TODO: Add rate limiting
  // TODO: Make this search paginated
  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonTitle>Admin Books</IonTitle>
        </IonToolbar>
      </IonHeader>

      <IonContent className="ion-padding">
        <div className="admin-books-header">
          <IonText>
            <h2>Admin Books Section</h2>
          </IonText>
          <IonButton
            shape="round"
            className="float-right-button"
            onClick={() => history.push('/admin/search-google-books')}
          >
            <IonIcon slot="icon-only" icon={addOutline}></IonIcon>
          </IonButton>
        </div>
        {loading ? <IonSpinner /> : <BookGrid books={books} />}
      </IonContent>
    </IonPage>
  );
}
