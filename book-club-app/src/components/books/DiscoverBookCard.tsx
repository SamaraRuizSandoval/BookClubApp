import { IonImg } from '@ionic/react';

import '../../styles/discover-book-card.css';
import { Book } from '../../types/book';

type DiscoverBookCardProps = {
  book: Book;
  onAddToCollection: (
    bookId: number,
    status: 'wishlist' | 'reading' | 'completed',
  ) => void;
};

export function DiscoverBooksCard({
  book,
  onAddToCollection,
}: DiscoverBookCardProps) {
  const imageUrl = book.book_images?.small_url;

  return (
    <div className="book-card" role="listitem" aria-label={book.title}>
      <div className="cover-wrap">
        <div className="cover-inner">
          {!imageUrl ? (
            <div className="cover-placeholder">
              <span className="ph-icon">📖</span>
              <span className="ph-title">
                {book.title || 'Cover not available'}
              </span>
            </div>
          ) : (
            <IonImg src={imageUrl} alt={book.title} />
          )}
        </div>

        <div className="cover-shine" aria-hidden="true"></div>

        {/* Hover actions */}
        <div className="card-actions">
          <div className="action-row">
            <button
              className="btn-collection reading"
              onClick={() => onAddToCollection(book.id, 'reading')}
            >
              📖 Reading
            </button>
            <button
              className="btn-collection want"
              onClick={() => onAddToCollection(book.id, 'wishlist')}
            >
              🔖 Want
            </button>
            <button
              className="btn-collection done"
              onClick={() => onAddToCollection(book.id, 'completed')}
            >
              ✓ Done
            </button>
          </div>
        </div>
      </div>

      <div className="book-meta">
        <p className="book-title">{book.title}</p>
        <p className="book-author">
          {book.authors?.length ? book.authors.join(', ') : 'Unknown Author'}
        </p>
      </div>
    </div>
  );
}
