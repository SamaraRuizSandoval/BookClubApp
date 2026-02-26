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
    <IonCard>
      <div className="flex-center">
        <IonImg
          className="book-cover"
          src={book.book_images?.thumbnail_url}
          alt={book.title}
        />
      </div>

      <IonCardHeader>
        <IonCardTitle className="book-card">{book.title}</IonCardTitle>
        <IonCardSubtitle> {book.authors.join(', ')} </IonCardSubtitle>
      </IonCardHeader>

      <IonCardContent>{book.publisher}</IonCardContent>
    </IonCard>
  );
}
