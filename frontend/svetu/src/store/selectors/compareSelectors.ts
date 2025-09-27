import { createSelector } from '@reduxjs/toolkit';
import { RootState } from '../index';

// Базовый селектор для получения всех items по категориям
const selectCompareItemsByCategory = (state: RootState) =>
  state.universalCompare?.itemsByCategory || {};

// Мемоизированный селектор для получения items конкретной категории
export const makeSelectCompareItemsByType = () =>
  createSelector(
    [selectCompareItemsByCategory, (_: RootState, type: string) => type],
    (itemsByCategory, type) => itemsByCategory[type] || []
  );

// Селектор для получения favorites
export const selectFavorites = (state: RootState) =>
  state.favorites?.items || [];
