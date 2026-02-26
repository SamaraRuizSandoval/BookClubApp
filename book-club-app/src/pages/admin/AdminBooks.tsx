import {
  IonContent,
  IonPage,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonText,
  IonSpinner,
} from '@ionic/react';
import { useEffect, useState } from 'react';

import api from '../../api/axios';
import { BookGrid } from '../../components/BooksGrid';
import { AuthTokenResponse } from '../../types/auth';
import { Book, BookResponse } from '../../types/book';

export function AdminBooks() {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(true);

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
  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonTitle>Admin Books</IonTitle>
        </IonToolbar>
      </IonHeader>

      <IonContent className="ion-padding">
        <IonText>
          <h2>Admin Books Section</h2>
          {loading ? <IonSpinner /> : <BookGrid books={books} />}
        </IonText>
      </IonContent>
    </IonPage>
  );
}
