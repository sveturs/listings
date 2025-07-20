-- Text search migration

CREATE TEXT SEARCH DICTIONARY public.unaccent_dict (
    TEMPLATE = public.unaccent,
    rules = 'unaccent' );



CREATE TEXT SEARCH CONFIGURATION public.english_unaccent (
    PARSER = pg_catalog."default" );
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR asciiword WITH english_stem;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR word WITH public.unaccent_dict, english_stem;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR numword WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR email WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR url WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR host WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR sfloat WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR version WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR hword_numpart WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR hword_part WITH public.unaccent_dict, english_stem;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR hword_asciipart WITH english_stem;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR numhword WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR asciihword WITH english_stem;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR hword WITH public.unaccent_dict, english_stem;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR url_path WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR file WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR "float" WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR "int" WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.english_unaccent
    ADD MAPPING FOR uint WITH simple;

CREATE TEXT SEARCH CONFIGURATION public.russian_unaccent (
    PARSER = pg_catalog."default" );
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR asciiword WITH english_stem;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR word WITH public.unaccent_dict, russian_stem;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR numword WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR email WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR url WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR host WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR sfloat WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR version WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR hword_numpart WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR hword_part WITH public.unaccent_dict, russian_stem;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR hword_asciipart WITH english_stem;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR numhword WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR asciihword WITH english_stem;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR hword WITH public.unaccent_dict, russian_stem;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR url_path WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR file WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR "float" WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR "int" WITH simple;
ALTER TEXT SEARCH CONFIGURATION public.russian_unaccent
    ADD MAPPING FOR uint WITH simple;