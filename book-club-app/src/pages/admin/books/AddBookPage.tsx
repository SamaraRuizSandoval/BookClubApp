import {
  IonContent,
  IonPage,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonItem,
  IonInput,
  IonTextarea,
  IonButton,
} from '@ionic/react';
import { useEffect, useState } from 'react';
import { useHistory, useLocation } from 'react-router-dom';

import api from '../../../api/axios';
import { Book, Chapter } from '../../../types/book';

type LocationState = {
  book: Book;
};

export function AddBookPage() {
  const location = useLocation<LocationState>();
  const history = useHistory();
  const book = location.state?.book;
  const [bookDraft, setBookDraft] = useState<Book | null>(null);
  const [chapters, setChapters] = useState<Chapter[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    if (book) {
      setBookDraft({
        ...book,
        chapters: [],
      });
    }
  }, [book]);

  const addChapter = () => {
    setChapters([
      ...chapters,
      { id: chapters.length, number: chapters.length + 1, title: '' },
    ]);
  };

  function updateChapter(
    index: number,
    field: keyof Chapter,
    value: string | number,
  ) {
    const updated = [...chapters];
    updated[index] = { ...updated[index], [field]: value };
    setChapters(updated);
  }

  const saveBook = async () => {
    try {
      setIsLoading(true);
      setErrorMessage(null); // clear previous error

      book.chapters = chapters;

      const token = localStorage.getItem('authToken');

      if (!token) {
        console.error('No auth token found');
        setIsLoading(false);
        return;
      }

      const response = await api.post<Book>(
        'books',
        {
          title: bookDraft?.title,
          authors: bookDraft?.authors,
          book_images: book.book_images,
          description: bookDraft?.description,
          page_count: bookDraft?.page_count,
          publisher: bookDraft?.publisher,
          published_date: bookDraft?.published_date,
          isbn_10: bookDraft?.isbn_10,
          isbn_13: bookDraft?.isbn_13,
          chapters: chapters,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        },
      );
      console.log('Book added:', response.data);
      setIsLoading(false);
      history.push('/admin/books');
    } catch (error: any) {
      setIsLoading(false);
      console.error('Error adding book:', error);
      setErrorMessage(
        error.response?.data?.error ||
          'An error occurred while adding the book',
      );
    }
  };

  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonTitle>Add Book</IonTitle>
        </IonToolbar>
      </IonHeader>
      <IonContent className="ion-padding">
        {!book ? (
          <p>Data not found</p>
        ) : (
          <>
            <h1>Add Book</h1>

            <img src={book.book_images?.small_url} alt={book.title} />
            {errorMessage && (
              <div className="error-message">{errorMessage}</div>
            )}
            <IonItem>
              <IonInput
                label="title"
                labelPlacement="stacked"
                placeholder="Enter book title"
                value={bookDraft?.title}
                onIonChange={(e) =>
                  setBookDraft({
                    ...bookDraft!,
                    title: e.detail.value || '',
                  })
                }
              ></IonInput>
            </IonItem>

            <IonItem>
              <IonInput
                label="Authors"
                labelPlacement="stacked"
                placeholder="Enter book authors, if multiple separate by comma"
                value={bookDraft?.authors?.join(', ')}
                onIonChange={(e) =>
                  setBookDraft({
                    ...bookDraft!,
                    authors:
                      e.detail.value?.split(', ').map((a) => a.trim()) || [],
                  })
                }
              ></IonInput>
            </IonItem>

            <IonItem>
              <IonInput
                label="Publisher"
                labelPlacement="stacked"
                placeholder="Enter book publisher"
                value={bookDraft?.publisher}
                onIonChange={(e) =>
                  setBookDraft({
                    ...bookDraft!,
                    publisher: e.detail.value || '',
                  })
                }
              ></IonInput>
            </IonItem>

            <IonItem>
              <IonInput
                label="Page Count"
                labelPlacement="stacked"
                placeholder="Enter book page count"
                value={bookDraft?.page_count}
                onIonChange={(e) =>
                  setBookDraft({
                    ...bookDraft!,
                    page_count: parseInt(e.detail.value || '0', 10),
                  })
                }
              ></IonInput>
            </IonItem>

            <IonItem>
              <IonTextarea
                className=""
                label="Description"
                labelPlacement="stacked"
                autoGrow={true}
                placeholder="Enter book description"
                value={bookDraft?.description}
                onIonChange={(e) =>
                  setBookDraft({
                    ...bookDraft!,
                    description: e.detail.value || '',
                  })
                }
              ></IonTextarea>
            </IonItem>

            <IonItem>
              <IonInput
                label="ISBN 10"
                labelPlacement="stacked"
                placeholder="Enter book ISBN 10"
                value={bookDraft?.isbn_10}
                onIonChange={(e) =>
                  setBookDraft({
                    ...bookDraft!,
                    isbn_10: e.detail.value || '',
                  })
                }
              ></IonInput>
            </IonItem>

            <IonItem>
              <IonInput
                label="ISBN 13"
                labelPlacement="stacked"
                placeholder="Enter book ISBN 13"
                value={bookDraft?.isbn_13}
                onIonChange={(e) =>
                  setBookDraft({
                    ...bookDraft!,
                    isbn_13: e.detail.value || '',
                  })
                }
              ></IonInput>
            </IonItem>

            <IonItem>
              <IonInput
                label="Published Date"
                labelPlacement="stacked"
                placeholder="Enter book published date"
                value={book.published_date}
              ></IonInput>
            </IonItem>

            <h2>Chapters</h2>
            {chapters.map((chapter) => (
              <IonItem key={chapter.id}>
                <IonInput
                  label="Chapter Number"
                  labelPlacement="stacked"
                  placeholder="Enter chapter number"
                  value={chapter.number}
                ></IonInput>
                <IonInput
                  label="Chapter Title"
                  labelPlacement="stacked"
                  placeholder="Enter chapter title"
                  value={chapter.title}
                  onIonChange={(e) =>
                    updateChapter(chapter.id, 'title', e.detail.value || '')
                  }
                ></IonInput>
              </IonItem>
            ))}

            <IonButton expand="block" onClick={addChapter}>
              Add Chapter
            </IonButton>

            <IonButton
              expand="block"
              className="ion-margin-top"
              onClick={saveBook}
              disabled={isLoading}
            >
              Save Book
            </IonButton>
          </>
        )}
      </IonContent>
    </IonPage>
  );
}
