import {
  IonContent,
  IonPage,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonText,
  IonSearchbar,
  IonSpinner,
} from '@ionic/react';
import { useEffect, useState } from 'react';

import api from '../../../api/axios';
import { BookGrid } from '../../../components/BooksGrid';
import { Book, BookResponse, BooksGoogleResponse } from '../../../types/book';

export function SearchGoogleBooks() {
  const [query, setQuery] = useState('');
  const [newBooks, setNewBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!query) return;

    const timeout = setTimeout(() => {
      console.log('Searching for books with query:', query);
      handleBookSearch(query);
    }, 500);

    return () => clearTimeout(timeout);
  }, [query]);

  const handleBookSearch = async (searchQuery: string) => {
    setLoading(true);

    const token = localStorage.getItem('authToken');

    if (!token) {
      console.error('No auth token found');
      setLoading(false);
      return;
    }

    try {
      const response = await api.get<BooksGoogleResponse>('/api/books', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        params: {
          q: searchQuery,
        },
      });

      // TODO: Check if the response is valid and handle errors

      setNewBooks(response.data.books);
      console.log(response.data.books);
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonTitle>Admin Books</IonTitle>
        </IonToolbar>
      </IonHeader>

      <IonContent className="ion-padding">
        <IonText>
          <h2>Search new books with Google Books API</h2>
        </IonText>
        <IonSearchbar
          debounce={1000}
          onIonInput={(e) => setQuery(e.detail.value!)}
        ></IonSearchbar>
        {loading ? <IonSpinner /> : <BookGrid books={newBooks} />}
      </IonContent>
    </IonPage>
  );
}
