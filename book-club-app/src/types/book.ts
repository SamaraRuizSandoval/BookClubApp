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
  authors: string[];
  description: string;
  book_images: BookImage;
  publisher: string;
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
