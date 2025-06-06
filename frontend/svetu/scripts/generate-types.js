#!/usr/bin/env node

/**
 * Script to generate TypeScript types from OpenAPI/Swagger specification
 * Usage: node scripts/generate-types.js
 */

const fs = require('fs');
const path = require('path');

// URL к swagger документации
const SWAGGER_URL = process.env.API_BASE_URL
  ? `${process.env.API_BASE_URL}/swagger/doc.json`
  : 'http://localhost:8080/swagger/doc.json';

// Путь для сохранения сгенерированных типов
const OUTPUT_PATH = path.join(__dirname, '../src/types/generated/openapi.ts');

async function fetchSwaggerSpec() {
  try {
    const response = await fetch(SWAGGER_URL);
    if (!response.ok) {
      throw new Error(`Failed to fetch swagger spec: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching swagger spec:', error);
    // Возвращаем пустую спецификацию если не удалось загрузить
    return { paths: {}, definitions: {} };
  }
}

function generateTypeFromSchema(schema, definitions) {
  if (!schema) return 'unknown';

  if (schema.$ref) {
    const refName = schema.$ref.split('/').pop();
    return refName;
  }

  if (schema.type === 'array') {
    const itemType = generateTypeFromSchema(schema.items, definitions);
    return `${itemType}[]`;
  }

  if (schema.type === 'object') {
    if (schema.properties) {
      const props = Object.entries(schema.properties)
        .map(([key, prop]) => {
          const required = schema.required?.includes(key);
          const type = generateTypeFromSchema(prop, definitions);
          return `  ${key}${required ? '' : '?'}: ${type};`;
        })
        .join('\n');
      return `{\n${props}\n}`;
    }
    return 'Record<string, unknown>';
  }

  switch (schema.type) {
    case 'string':
      return schema.enum
        ? schema.enum.map((v) => `'${v}'`).join(' | ')
        : 'string';
    case 'integer':
    case 'number':
      return 'number';
    case 'boolean':
      return 'boolean';
    default:
      return 'unknown';
  }
}

function generateInterface(name, schema, definitions) {
  if (!schema.properties) return '';

  const props = Object.entries(schema.properties)
    .map(([key, prop]) => {
      const required = schema.required?.includes(key);
      const type = generateTypeFromSchema(prop, definitions);
      return `  ${key}${required ? '' : '?'}: ${type};`;
    })
    .join('\n');

  return `export interface ${name} {\n${props}\n}`;
}

async function generateTypes() {
  console.log('Fetching swagger specification...');
  const spec = await fetchSwaggerSpec();

  console.log('Generating TypeScript types...');

  let output = `/**
 * This file is auto-generated from OpenAPI specification
 * Generated at: ${new Date().toISOString()}
 * Do not edit manually
 */

`;

  // Генерируем интерфейсы из definitions
  if (spec.definitions) {
    Object.entries(spec.definitions).forEach(([name, schema]) => {
      const interfaceCode = generateInterface(name, schema, spec.definitions);
      if (interfaceCode) {
        output += interfaceCode + '\n\n';
      }
    });
  }

  // Создаем директорию если не существует
  const dir = path.dirname(OUTPUT_PATH);
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }

  // Сохраняем файл
  fs.writeFileSync(OUTPUT_PATH, output);
  console.log(`Types generated successfully at: ${OUTPUT_PATH}`);
}

// Запускаем генерацию
generateTypes().catch(console.error);
