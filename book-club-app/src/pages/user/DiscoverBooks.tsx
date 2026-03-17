import { IonContent, IonSpinner } from '@ionic/react';

import { api } from '../../api/apiClient';
import { addBookToUserCollection } from '../../api/userBooksApi';
import { useAuth } from '../../context/AuthContext';
import '../../styles/discover-books.css';

import { DiscoverBooksCard } from '../../components/books/DiscoverBookCard';

import { useEffect, useState } from 'react';

import { Book, BookResponse } from '../../types/book';

export function DiscoverBooks() {
  const { auth } = useAuth();
  const user = auth.user;

  const [books, setBooks] = useState<Book[]>([]);
  const [loadingBooks, setLoadingBooks] = useState(true);

  useEffect(() => {
    async function fetchBooks() {
      try {
        const response = await api.get<BookResponse>('/books');
        setBooks(response.data.books);
      } catch (error) {
        console.error(error);
      } finally {
        setLoadingBooks(false);
      }
    }
    fetchBooks();
  }, []);

  async function handleAddToCollection(
    bookId: number,
    status: 'wishlist' | 'reading' | 'completed',
  ) {
    if (!user) return;

    try {
      await addBookToUserCollection(user.id, bookId, status);

      console.log('Book added!');
    } catch (error) {
      console.error('Failed to add book', error);
    }
  }

  return (
    <>
      <IonContent className="body-bg">
        <style>
          @import
          url('https://fonts.googleapis.com/css2?family=DM+Serif+Display:ital@0;1&display=swap');
        </style>

        {/* -- GREETING -- */}
        <div className="padding-top padding-side">
          <div className="greeting-hero">
            <div className="greeting-text">
              <h2>
                Happy reading,
                <br />
                <em>{user?.username}.</em>
              </h2>
              <p>
                Explore the full catalogue and add books to your collections.
                Hover a cover to get started.
              </p>
            </div>
            <div className="greeting-stats" aria-label="Reading stats">
              <div className="g-stat">
                <span className="g-stat-val" id="statReading">
                  0
                </span>
                <span className="g-stat-lbl">Reading</span>
              </div>
              <div className="g-stat-divider"></div>
              <div className="g-stat">
                <span className="g-stat-val" id="statWant">
                  0
                </span>
                <span className="g-stat-lbl">Want to Read</span>
              </div>
              <div className="g-stat-divider"></div>
              <div className="g-stat">
                <span className="g-stat-val" id="statDone">
                  0
                </span>
                <span className="g-stat-lbl">Finished</span>
              </div>
            </div>
          </div>
        </div>

        {/* -- SECTION HEADER -- */}
        <div className="padding-side">
          <div className="section-header">
            <div className="section-left">
              <div className="section-eyebrow">✦ Catalogue</div>
              <h1 className="section-title">All Available Books</h1>
            </div>
            <span id="bookCount"></span>
          </div>
        </div>

        <div className="padding-side">
          <div
            className="books-grid"
            id="booksGrid"
            role="list"
            aria-label="Books catalogue"
          >
            {loadingBooks ? (
              <IonSpinner />
            ) : (
              <>
                {books.map((book) => (
                  <DiscoverBooksCard
                    key={book.id}
                    book={book}
                    onAddToCollection={handleAddToCollection}
                  />
                ))}
              </>
            )}
          </div>
        </div>
      </IonContent>
    </>
  );
}
