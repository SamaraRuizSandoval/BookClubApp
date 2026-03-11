import { IonCol, IonGrid, IonRow } from '@ionic/react';

import { Book } from '../types/book';

import { BookCard } from './BookCard';

type BookGridProps = {
  books: Book[];
};
export function BookGrid({ books }: BookGridProps) {
  return (
    <IonGrid>
      <IonRow>
        {books.map((book) => (
          <IonCol size="8" sizeMd="4" sizeLg="3" key={book.id}>
            <BookCard book={book} />
          </IonCol>
        ))}
      </IonRow>
    </IonGrid>
  );
}
