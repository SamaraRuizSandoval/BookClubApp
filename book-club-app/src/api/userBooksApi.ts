import { api } from './apiClient';

export type BookStatus = 'wishlist' | 'reading' | 'completed';

export async function addBookToUserCollection(
  userId: number,
  bookId: number,
  status: BookStatus,
) {
  const response = await api.post(`/users/${userId}/books`, null, {
    params: {
      book_id: bookId,
      status: status,
    },
  });

  return response.data;
}
