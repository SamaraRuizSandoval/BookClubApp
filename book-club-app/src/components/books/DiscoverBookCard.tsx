import { IonImg } from '@ionic/react';

import { useToast } from '../../context/ToastContext';

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
  const smallUrl = book.book_images?.small_url;
  const mediumUrl = book.book_images?.medium_url;
  const { show } = useToast();

  const imageUrl = mediumUrl || smallUrl;

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
              onClick={() => {
                onAddToCollection(book.id, 'reading');
                show(`"${book.title}"📖 Added to Reading`, 'primary');
              }}
            >
              📖 Reading
            </button>
            <button
              className="btn-collection want"
              onClick={() => {
                onAddToCollection(book.id, 'wishlist');
                show(`"${book.title}"🔖 Added to Want to Read`, 'primary');
              }}
            >
              🔖 Want
            </button>
            <button
              className="btn-collection done"
              onClick={() => {
                onAddToCollection(book.id, 'completed');
                show(`"${book.title}"✓ Marked as Read`, 'primary');
              }}
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
