import {
  IonCard,
  IonCardHeader,
  IonCardTitle,
  IonCardSubtitle,
  IonCardContent,
  IonImg,
} from '@ionic/react';

import { Book } from '../types/book';

type BookCardProps = {
  book: Book;
};

export function BookCard({ book }: BookCardProps) {
  return (
    <IonCard className="book-card">
      <div className="flex-center">
        <IonImg
          className="book-cover"
          src={
            book.book_images?.small_url ||
            'https://cdnattic.atticbooks.co.ke/img/R741540.jpg'
          }
          alt={book.title}
        />
      </div>

      <IonCardHeader className="book-card">
        <IonCardTitle className="book-card-title">{book.title}</IonCardTitle>
        <IonCardSubtitle>
          {' '}
          {book.authors?.length
            ? book.authors.join(', ')
            : 'Unknown Author'}{' '}
        </IonCardSubtitle>
      </IonCardHeader>

      <IonCardContent>{book.publisher}</IonCardContent>
    </IonCard>
  );
}
