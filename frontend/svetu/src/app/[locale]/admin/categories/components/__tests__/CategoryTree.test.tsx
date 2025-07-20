import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CategoryTree from '../CategoryTree';
import { Category } from '@/services/admin';

// Mock the adminApi
jest.mock('@/services/admin', () => ({
  adminApi: {
    updateFieldTranslation: jest.fn(),
  },
}));

// Mock the TranslationStatus component
jest.mock('@/components/attributes/TranslationStatus', () => ({
  TranslationStatus: ({ compact }: { compact?: boolean }) => (
    <div data-testid="translation-status">
      TranslationStatus {compact ? 'compact' : 'full'}
    </div>
  ),
}));

// Mock the InlineTranslationEditor component
jest.mock('@/components/attributes/InlineTranslationEditor', () => ({
  InlineTranslationEditor: ({
    translations,
    onSave: _onSave,
  }: {
    translations: Record<string, string>;
    onSave: (translations: Record<string, string>) => void;
  }) => <div data-testid="inline-translation-editor">{translations.en}</div>,
}));

const mockCategories: Category[] = [
  {
    id: 1,
    name: 'Electronics',
    slug: 'electronics',
    parent_id: undefined,
    is_active: true,
    icon: '📱',
    items_count: 50,
    created_at: '2024-01-01',
    updated_at: '2024-01-01',
  },
  {
    id: 2,
    name: 'Smartphones',
    slug: 'smartphones',
    parent_id: 1,
    is_active: true,
    icon: '📱',
    items_count: 30,
    created_at: '2024-01-01',
    updated_at: '2024-01-01',
  },
  {
    id: 3,
    name: 'Laptops',
    slug: 'laptops',
    parent_id: 1,
    is_active: false,
    icon: '💻',
    items_count: 20,
    created_at: '2024-01-01',
    updated_at: '2024-01-01',
  },
  {
    id: 4,
    name: 'Clothing',
    slug: 'clothing',
    parent_id: undefined,
    is_active: true,
    icon: '👕',
    items_count: 100,
    created_at: '2024-01-01',
    updated_at: '2024-01-01',
  },
];

describe('CategoryTree', () => {
  const mockOnEdit = jest.fn();
  const mockOnDelete = jest.fn();
  const mockOnManageAttributes = jest.fn();
  const mockOnReorder = jest.fn();
  const mockOnMove = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders categories correctly', () => {
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onManageAttributes={mockOnManageAttributes}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    expect(screen.getByText('Electronics')).toBeInTheDocument();
    expect(screen.getByText('Clothing')).toBeInTheDocument();
  });

  it('displays category icons when available', () => {
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    // Multiple phone icons may exist (parent and child categories)
    const phoneIcons = screen.getAllByText('📱');
    expect(phoneIcons.length).toBeGreaterThan(0);
    expect(screen.getByText('👕')).toBeInTheDocument();
  });

  it('shows inactive badge for inactive categories', () => {
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    // Laptops is inactive
    const inactiveBadges = screen.getAllByText('common.inactive');
    expect(inactiveBadges).toHaveLength(1);
  });

  it('displays items count when available', () => {
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    expect(screen.getByText('50')).toBeInTheDocument(); // Electronics
    expect(screen.getByText('100')).toBeInTheDocument(); // Clothing
  });

  it('expands and collapses child categories', async () => {
    const user = userEvent.setup();
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    // Initially, child categories should be visible
    expect(screen.getByText('Smartphones')).toBeInTheDocument();
    expect(screen.getByText('Laptops')).toBeInTheDocument();

    // Find and click the expand/collapse button for Electronics
    const expandButtons = screen.getAllByLabelText('common.collapse');
    await user.click(expandButtons[0]);

    // Child categories should be hidden
    await waitFor(() => {
      expect(screen.queryByText('Smartphones')).not.toBeInTheDocument();
      expect(screen.queryByText('Laptops')).not.toBeInTheDocument();
    });

    // Click again to expand
    const collapseButtons = screen.getAllByLabelText('common.expand');
    await user.click(collapseButtons[0]);

    // Child categories should be visible again
    await waitFor(() => {
      expect(screen.getByText('Smartphones')).toBeInTheDocument();
      expect(screen.getByText('Laptops')).toBeInTheDocument();
    });
  });

  it('calls onEdit when edit action is clicked', async () => {
    const user = userEvent.setup();
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    // Open dropdown for Electronics (first dropdown button)
    const dropdownButtons = screen.getAllByLabelText('');
    await user.click(dropdownButtons[0]);

    // Click edit
    const editButtons = screen.getAllByText('common.edit');
    await user.click(editButtons[0]);

    expect(mockOnEdit).toHaveBeenCalledWith(mockCategories[0]);
  });

  it('calls onDelete when delete action is clicked', async () => {
    const user = userEvent.setup();
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    // Open dropdown for Electronics (first dropdown button)
    const dropdownButtons = screen.getAllByLabelText('');
    await user.click(dropdownButtons[0]);

    // Click delete
    const deleteButtons = screen.getAllByText('common.delete');
    await user.click(deleteButtons[0]);

    expect(mockOnDelete).toHaveBeenCalledWith(mockCategories[0]);
  });

  it('calls onManageAttributes when attributes action is clicked', async () => {
    const user = userEvent.setup();
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onManageAttributes={mockOnManageAttributes}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    // Open dropdown for Electronics (first dropdown button)
    const dropdownButtons = screen.getAllByLabelText('');
    await user.click(dropdownButtons[0]);

    // Click attributes
    const attributesButton = screen.getByText('sections.attributes');
    await user.click(attributesButton);

    expect(mockOnManageAttributes).toHaveBeenCalledWith(mockCategories[0]);
  });

  it('toggles category active status', async () => {
    const user = userEvent.setup();
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    // Open dropdown for active category (Electronics - first dropdown button)
    const dropdownButtons = screen.getAllByLabelText('');
    await user.click(dropdownButtons[0]);

    // Click deactivate
    const deactivateButton = screen.getByText('common.deactivate');
    await user.click(deactivateButton);

    expect(mockOnEdit).toHaveBeenCalledWith({
      ...mockCategories[0],
      is_active: false,
    });
  });

  it('filters categories based on active status', async () => {
    const user = userEvent.setup();
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    // Initially, all categories should be visible
    expect(screen.getByText('Laptops')).toBeInTheDocument();

    // Uncheck "Show Inactive"
    const showInactiveCheckbox = screen.getByRole('checkbox');
    await user.click(showInactiveCheckbox);

    // Inactive category should be hidden
    await waitFor(() => {
      expect(screen.queryByText('Laptops')).not.toBeInTheDocument();
    });

    // Active categories should still be visible
    expect(screen.getByText('Electronics')).toBeInTheDocument();
    expect(screen.getByText('Smartphones')).toBeInTheDocument();
  });

  it('displays no data message when categories array is empty', () => {
    render(
      <CategoryTree
        categories={[]}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    expect(screen.getByText('common.noData')).toBeInTheDocument();
  });

  it('renders hierarchical lines correctly', () => {
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    // Check that child categories have hierarchy indicators
    const categoryNodes = document.querySelectorAll('.category-node');
    expect(categoryNodes.length).toBeGreaterThan(0);
  });

  it('displays active/inactive status indicators', () => {
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    // Check for status indicators in the filter legend
    expect(screen.getByText('common.active')).toBeInTheDocument();
    const inactiveLabels = screen.getAllByText('common.inactive');
    expect(inactiveLabels.length).toBeGreaterThan(0);
  });

  it('renders translation status component for each category', () => {
    render(
      <CategoryTree
        categories={mockCategories}
        onEdit={mockOnEdit}
        onDelete={mockOnDelete}
        onReorder={mockOnReorder}
        onMove={mockOnMove}
      />
    );

    const translationStatuses = screen.getAllByTestId('translation-status');
    expect(translationStatuses.length).toBe(mockCategories.length);
  });
});
