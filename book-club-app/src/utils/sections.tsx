import React from 'react';

import { BooksSectionLayout } from './../components/BookSectionLayout';

export const ReadingSection = () => (
  <BooksSectionLayout sectionKey="reading" sectionTitle="Reading" />
);

export const WishlistSection = () => (
  <BooksSectionLayout sectionKey="wishlist" sectionTitle="Wishlist" />
);

export const CompletedSection = () => (
  <BooksSectionLayout sectionKey="completed" sectionTitle="Completed" />
);
