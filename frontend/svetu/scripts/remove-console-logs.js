#!/usr/bin/env node

/**
 * Script to remove console.* statements from production code
 * Keeps console.error and console.warn for error tracking
 */

const fs = require('fs');
const path = require('path');

// Directories to process
const srcDir = path.join(__dirname, '..', 'src');

// Patterns to find TypeScript/JavaScript files
const patterns = ['**/*.ts', '**/*.tsx', '**/*.js', '**/*.jsx'];

// Files or directories to exclude
const excludePatterns = [
  '**/node_modules/**',
  '**/*.test.*',
  '**/*.spec.*',
  '**/test/**',
  '**/tests/**',
  '**/__tests__/**',
  '**/debug-*/**', // Debug pages
  '**/scripts/**',
];

// Console methods to remove (keep error and warn for production debugging)
const consoleMethodsToRemove = [
  'log',
  'info',
  'debug',
  'trace',
  'time',
  'timeEnd',
  'group',
  'groupEnd',
];

// Regex pattern to match console statements
const createConsoleRegex = () => {
  const methods = consoleMethodsToRemove.join('|');
  // Match console.method(...) including multiline
  return new RegExp(`console\\.(${methods})\\([^)]*\\)(?:\\s*;)?`, 'g');
};

// More complex regex for multiline console statements
const createMultilineConsoleRegex = () => {
  const methods = consoleMethodsToRemove.join('|');
  return new RegExp(
    `console\\.(${methods})\\([^)]*(?:\\n[^)]*)*\\)(?:\\s*;)?`,
    'gm'
  );
};

let totalRemoved = 0;
let filesProcessed = 0;

// Process a single file
function processFile(filePath) {
  try {
    let content = fs.readFileSync(filePath, 'utf8');
    const originalContent = content;

    // Count occurrences before removal
    const simpleRegex = createConsoleRegex();
    const multilineRegex = createMultilineConsoleRegex();

    const matches = [
      ...(content.match(simpleRegex) || []),
      ...(content.match(multilineRegex) || []),
    ];

    if (matches.length > 0) {
      // Remove simple console statements
      content = content.replace(simpleRegex, '');

      // Remove multiline console statements
      content = content.replace(multilineRegex, '');

      // Clean up empty lines left behind
      content = content.replace(/^\s*[\r\n]/gm, '\n');
      content = content.replace(/\n\n\n+/g, '\n\n');

      // Only write if content changed
      if (content !== originalContent) {
        fs.writeFileSync(filePath, content, 'utf8');
        console.log(
          `✓ ${path.relative(srcDir, filePath)}: Removed ${matches.length} console.${consoleMethodsToRemove.join('/console.')} statements`
        );
        totalRemoved += matches.length;
        filesProcessed++;
      }
    }
  } catch (error) {
    console.error(`✗ Error processing ${filePath}:`, error.message);
  }
}

// Recursively find files
function findFiles(dir, fileList = []) {
  const files = fs.readdirSync(dir);

  files.forEach((file) => {
    const filePath = path.join(dir, file);
    const stat = fs.statSync(filePath);

    // Skip excluded directories
    const relativePath = path.relative(srcDir, filePath);
    if (
      excludePatterns.some((pattern) => {
        const regex = new RegExp(
          pattern.replace(/\*/g, '.*').replace(/\//g, '\\/')
        );
        return regex.test(relativePath);
      })
    ) {
      return;
    }

    if (stat.isDirectory()) {
      findFiles(filePath, fileList);
    } else if (/\.(ts|tsx|js|jsx)$/.test(file)) {
      fileList.push(filePath);
    }
  });

  return fileList;
}

// Main function
function main() {
  console.log('🧹 Removing console statements from production code...\n');

  const files = findFiles(srcDir);
  files.forEach(processFile);

  console.log(
    `\n✅ Complete! Removed ${totalRemoved} console statements from ${filesProcessed} files.`
  );

  if (totalRemoved > 0) {
    console.log('\n⚠️  Remember to:');
    console.log('  1. Review the changes with git diff');
    console.log('  2. Run tests to ensure nothing broke');
    console.log(
      '  3. Consider using a logger service for production debugging'
    );
  }
}

// Run the script
main();
