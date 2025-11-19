/* Do not change, this code is generated from Golang structs */

export interface Chapter {
  id: number;
  number: number;
  title: string;
}
export interface BookImages {
  thumbnail_url?: string;
  small_url?: string;
  medium_url?: string;
  large_url?: string;
}
export interface JSONDate {}
export interface Book {
  id: number;
  title: string;
  authors: string[];
  publisher: string;
  published_date: JSONDate;
  description?: string;
  page_count?: number;
  isbn_13: string;
  isbn_10?: string;
  book_images: BookImages;
  chapters: Chapter[];
}
