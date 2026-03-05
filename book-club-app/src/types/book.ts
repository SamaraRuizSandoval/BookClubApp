export type Chapter = {
  id: number;
  number: number;
  title: string;
};

export type BookImage = {
  large_url: string;
  medium_url: string;
  small_url: string;
  thumbnail_url: string;
};

export type Book = {
  id: number;
  title: string;
  authors: string[] | null;
  description: string;
  book_images: BookImage;
  publisher: string | null;
  published_date: string;
  isbn_10: string;
  isbn_13: string;
  chapters: Chapter[];
};

export type BookResponse = {
  books: Book[];
  limit: number;
  page: number;
  total_items: number;
  total_pages: number;
};

export type BooksGoogleResponse = {
  books: Book[];
};
