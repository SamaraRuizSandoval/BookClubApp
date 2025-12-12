export type BookStatus = 'reading' | 'wishlist' | 'completed';

export interface Book {
  id: string;
  title: string;
  status: BookStatus;
}

export const allBooks: Book[] = [
  { id: 'hp1', title: 'Harry Potter 1', status: 'reading' },
  { id: 'hp2', title: 'Harry Potter 2', status: 'reading' },
  { id: 'hp3', title: 'Harry Potter 3', status: 'wishlist' },
  { id: 'acotar', title: 'A Court of Thorns and Roses', status: 'completed' },
];
