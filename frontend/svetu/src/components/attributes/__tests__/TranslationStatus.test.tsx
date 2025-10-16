import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { TranslationStatus } from '../TranslationStatus';
import { adminApi } from '@/services/admin';

// Mock the adminApi
jest.mock('@/services/admin', () => ({
  adminApi: {
    getTranslationStatus: jest.fn(),
  },
}));

const mockTranslationStatus = {
  entity_type: 'category',
  entity_id: 1,
  languages: {
    en: {
      is_translated: true,
      is_machine_translated: false,
      is_verified: true,
      translated_at: '2024-01-01',
    },
    ru: {
      is_translated: true,
      is_machine_translated: true,
      is_verified: false,
      translated_at: '2024-01-01',
    },
    sr: {
      is_translated: false,
      is_machine_translated: false,
      is_verified: false,
      translated_at: null,
    },
  },
};

describe('TranslationStatus', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('displays loading state initially', () => {
    (adminApi.getTranslationStatus as jest.Mock).mockImplementation(
      () => new Promise(() => {}) // Never resolves to keep loading state
    );

    render(<TranslationStatus entityType="category" entityId={1} />);

    expect(screen.getByText('loading')).toBeInTheDocument();
    expect(document.querySelector('.loading-spinner')).toBeInTheDocument();
  });

  it('fetches translation status on mount', async () => {
    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([
      mockTranslationStatus,
    ]);

    render(<TranslationStatus entityType="category" entityId={1} />);

    await waitFor(() => {
      expect(adminApi.getTranslationStatus).toHaveBeenCalledWith('category', [
        1,
      ]);
    });
  });

  it('displays full status view by default', async () => {
    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([
      mockTranslationStatus,
    ]);

    render(<TranslationStatus entityType="category" entityId={1} />);

    await waitFor(() => {
      expect(screen.getByText('translations.status')).toBeInTheDocument();
      expect(screen.getByText('English')).toBeInTheDocument();
      expect(screen.getByText('Русский')).toBeInTheDocument();
      expect(screen.getByText('Српски')).toBeInTheDocument();
    });
  });

  it('displays compact view when compact prop is true', async () => {
    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([
      mockTranslationStatus,
    ]);

    render(
      <TranslationStatus entityType="category" entityId={1} compact={true} />
    );

    await waitFor(() => {
      // In compact mode, only icons are shown
      expect(screen.getByText('✅')).toBeInTheDocument(); // English verified
      expect(screen.getByText('🤖')).toBeInTheDocument(); // Russian machine translated
      expect(screen.getByText('❌')).toBeInTheDocument(); // Serbian not translated
    });
  });

  it('shows correct icons for different translation states', async () => {
    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([
      mockTranslationStatus,
    ]);

    render(
      <TranslationStatus entityType="category" entityId={1} compact={true} />
    );

    await waitFor(() => {
      const verifiedIcon = screen.getByText('✅');
      const machineIcon = screen.getByText('🤖');
      const notTranslatedIcon = screen.getByText('❌');

      expect(verifiedIcon).toBeInTheDocument();
      expect(machineIcon).toBeInTheDocument();
      expect(notTranslatedIcon).toBeInTheDocument();
    });
  });

  it('displays translate button when not all languages are translated', async () => {
    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([
      mockTranslationStatus,
    ]);
    const mockOnTranslateClick = jest.fn();

    render(
      <TranslationStatus
        entityType="category"
        entityId={1}
        onTranslateClick={mockOnTranslateClick}
        compact={true}
      />
    );

    await waitFor(() => {
      const translateButton = screen.getByTitle('translations.translate');
      expect(translateButton).toBeInTheDocument();
      expect(translateButton).toHaveTextContent('🌍');
    });
  });

  it('calls onTranslateClick when translate button is clicked', async () => {
    const user = userEvent.setup();
    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([
      mockTranslationStatus,
    ]);
    const mockOnTranslateClick = jest.fn();

    render(
      <TranslationStatus
        entityType="category"
        entityId={1}
        onTranslateClick={mockOnTranslateClick}
        compact={true}
      />
    );

    await waitFor(() => {
      const translateButton = screen.getByTitle('translations.translate');
      expect(translateButton).toBeInTheDocument();
    });

    const translateButton = screen.getByTitle('translations.translate');
    await user.click(translateButton);

    expect(mockOnTranslateClick).toHaveBeenCalled();
  });

  it('does not show translate button when all languages are translated', async () => {
    const allTranslatedStatus = {
      ...mockTranslationStatus,
      languages: {
        en: {
          is_translated: true,
          is_machine_translated: false,
          is_verified: true,
          translated_at: '2024-01-01',
        },
        ru: {
          is_translated: true,
          is_machine_translated: false,
          is_verified: true,
          translated_at: '2024-01-01',
        },
        sr: {
          is_translated: true,
          is_machine_translated: false,
          is_verified: true,
          translated_at: '2024-01-01',
        },
      },
    };

    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([
      allTranslatedStatus,
    ]);

    render(
      <TranslationStatus
        entityType="category"
        entityId={1}
        onTranslateClick={jest.fn()}
      />
    );

    await waitFor(() => {
      expect(screen.queryByTitle('translations.translate')).not.toBeInTheDocument();
      expect(screen.getByText('translations.allTranslated')).toBeInTheDocument();
    });
  });

  it('displays manual translation icon for non-verified manual translations', async () => {
    const manualTranslationStatus = {
      ...mockTranslationStatus,
      languages: {
        en: {
          is_translated: true,
          is_machine_translated: false,
          is_verified: false,
          translated_at: '2024-01-01',
        },
        ru: mockTranslationStatus.languages.ru,
        sr: mockTranslationStatus.languages.sr,
      },
    };

    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([
      manualTranslationStatus,
    ]);

    render(
      <TranslationStatus entityType="category" entityId={1} compact={true} />
    );

    await waitFor(() => {
      expect(screen.getByText('✏️')).toBeInTheDocument(); // Manual translation icon
    });
  });

  it('shows tooltips with correct language status in compact mode', async () => {
    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([
      mockTranslationStatus,
    ]);

    render(
      <TranslationStatus entityType="category" entityId={1} compact={true} />
    );

    await waitFor(() => {
      const tooltips = document.querySelectorAll('.tooltip');
      expect(tooltips.length).toBe(3); // One for each language
    });
  });

  it('handles API error gracefully', async () => {
    (adminApi.getTranslationStatus as jest.Mock).mockRejectedValue(
      new Error('API Error')
    );
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation();

    render(<TranslationStatus entityType="category" entityId={1} />);

    await waitFor(() => {
      expect(consoleSpy).toHaveBeenCalledWith(
        'Failed to fetch translation status:',
        expect.any(Error)
      );
    });

    consoleSpy.mockRestore();
  });

  it('returns null when no status data is available', async () => {
    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([]);

    const { container } = render(
      <TranslationStatus entityType="category" entityId={1} />
    );

    await waitFor(() => {
      expect(container.firstChild).toBeNull();
    });
  });

  it('re-fetches data when entityId changes', async () => {
    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([
      mockTranslationStatus,
    ]);

    const { rerender } = render(
      <TranslationStatus entityType="category" entityId={1} />
    );

    await waitFor(() => {
      expect(adminApi.getTranslationStatus).toHaveBeenCalledWith('category', [
        1,
      ]);
    });

    rerender(<TranslationStatus entityType="category" entityId={2} />);

    await waitFor(() => {
      expect(adminApi.getTranslationStatus).toHaveBeenCalledWith('category', [
        2,
      ]);
    });
  });

  it('applies correct color classes for different states', async () => {
    (adminApi.getTranslationStatus as jest.Mock).mockResolvedValue([
      mockTranslationStatus,
    ]);

    render(
      <TranslationStatus entityType="category" entityId={1} compact={true} />
    );

    await waitFor(() => {
      const tooltipDivs = document.querySelectorAll('.tooltip');

      expect(tooltipDivs[0]).toHaveClass('text-success'); // English verified
      expect(tooltipDivs[1]).toHaveClass('text-warning'); // Russian machine translated
      expect(tooltipDivs[2]).toHaveClass('text-error'); // Serbian not translated
    });
  });
});
