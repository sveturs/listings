import { NextRequest, NextResponse } from 'next/server';
import { exec } from 'child_process';
import { promisify } from 'util';
import fs from 'fs';
import path from 'path';

const execAsync = promisify(exec);

// Путь к файлу с районами
const DISTRICTS_FILE = path.join(
  process.cwd(),
  'src/app/[locale]/examples/novi-sad-districts/page.tsx'
);

// GET - получить список текущих районов
export async function GET() {
  try {
    const content = fs.readFileSync(DISTRICTS_FILE, 'utf8');

    // Извлекаем районы из файла
    const districtsMatch = content.match(/const districts = \[([\s\S]*?)\];/);
    if (!districtsMatch) {
      return NextResponse.json({ districts: [] });
    }

    // Парсим id и названия районов
    const idMatches = [...content.matchAll(/id:\s*['"`]([^'"`]+)['"`]/g)];
    const nameMatches = [...content.matchAll(/name:\s*['"`]([^'"`]+)['"`]/g)];
    const nameEnMatches = [
      ...content.matchAll(/nameEn:\s*['"`]([^'"`]+)['"`]/g),
    ];

    const districts = idMatches.map((match, index) => ({
      id: match[1],
      name: nameMatches[index]?.[1] || 'Unknown',
      nameEn: nameEnMatches[index]?.[1] || match[1],
      points: 0, // Будем считать точки позже если нужно
    }));

    return NextResponse.json({ districts });
  } catch (error) {
    console.error('Error reading districts:', error);
    return NextResponse.json({ districts: [] }, { status: 500 });
  }
}

// DELETE - удалить районы
export async function DELETE(request: NextRequest) {
  try {
    const { districtIds } = await request.json();

    if (
      !districtIds ||
      !Array.isArray(districtIds) ||
      districtIds.length === 0
    ) {
      return NextResponse.json(
        { error: 'No districts specified' },
        { status: 400 }
      );
    }

    // Выполняем скрипт удаления
    const scriptPath = path.join(
      process.cwd(),
      'scripts/novi-sad-districts/safe_remove_districts.js'
    );
    const command = `node ${scriptPath} ${districtIds.join(' ')}`;

    const { stdout, stderr } = await execAsync(command);

    if (stderr && !stderr.includes('✅')) {
      throw new Error(stderr);
    }

    // Парсим результат из stdout
    const removedMatch = stdout.match(/✅ Успешно удалено: (\d+) районов/);
    const remainingMatch = stdout.match(/📍 Осталось районов: (\d+)/);

    return NextResponse.json({
      success: true,
      removed: removedMatch ? parseInt(removedMatch[1]) : 0,
      remaining: remainingMatch ? parseInt(remainingMatch[1]) : 0,
      message: `Успешно удалено ${removedMatch ? removedMatch[1] : 0} районов`,
      output: stdout,
    });
  } catch (error) {
    console.error('Error deleting districts:', error);
    return NextResponse.json(
      {
        error: 'Failed to delete districts',
        details: error instanceof Error ? error.message : String(error),
      },
      { status: 500 }
    );
  }
}

// POST - добавить новый район
export async function POST(request: NextRequest) {
  try {
    const {
      name,
      city = 'Novi Sad',
      country = 'Serbia',
      source = 'nominatim',
      osmId,
      osmType,
    } = await request.json();

    if (!name) {
      return NextResponse.json(
        { error: 'District name is required' },
        { status: 400 }
      );
    }

    // Выбираем скрипт в зависимости от источника
    let scriptPath, command;

    if (source === 'overpass' && osmId && osmType) {
      // Используем Overpass API скрипт
      scriptPath = path.join(
        process.cwd(),
        'scripts/novi-sad-districts/add_district_overpass.js'
      );
      command = `node ${scriptPath} "${name}" ${osmId} ${osmType}`;
    } else {
      // Используем обычный Nominatim скрипт
      scriptPath = path.join(
        process.cwd(),
        'scripts/novi-sad-districts/auto_add_district.js'
      );
      command = `node ${scriptPath} "${name}" "${city}" "${country}"`;
    }

    const { stdout, stderr } = await execAsync(command);

    if (stderr || stdout.includes('❌')) {
      throw new Error(stderr || stdout);
    }

    // Проверяем успешность
    const successMatch = stdout.includes('🎉 Готово!');

    return NextResponse.json({
      success: successMatch,
      message: successMatch
        ? `Район "${name}" успешно добавлен`
        : 'Не удалось добавить район',
      output: stdout,
    });
  } catch (error) {
    console.error('Error adding district:', error);
    return NextResponse.json(
      {
        error: 'Failed to add district',
        details: error instanceof Error ? error.message : String(error),
      },
      { status: 500 }
    );
  }
}
