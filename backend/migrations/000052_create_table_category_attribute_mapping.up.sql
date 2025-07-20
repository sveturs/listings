-- Migration for table: category_attribute_mapping

CREATE TABLE public.category_attribute_mapping (
    category_id integer NOT NULL,
    attribute_id integer NOT NULL,
    is_enabled boolean DEFAULT true,
    is_required boolean DEFAULT false,
    sort_order integer DEFAULT 0,
    custom_component character varying(255),
    show_in_card boolean,
    show_in_list boolean
);

CREATE INDEX idx_category_attr_mapping ON public.category_attribute_mapping USING btree (category_id, is_enabled) WHERE (is_enabled = true);

CREATE INDEX idx_category_attribute_map_attr_id ON public.category_attribute_mapping USING btree (attribute_id);

CREATE INDEX idx_category_attribute_map_cat_id ON public.category_attribute_mapping USING btree (category_id);

CREATE INDEX idx_category_attribute_mapping_custom_component ON public.category_attribute_mapping USING btree (custom_component);

ALTER TABLE ONLY public.category_attribute_mapping
    ADD CONSTRAINT category_attribute_mapping_pkey PRIMARY KEY (category_id, attribute_id);

ALTER TABLE ONLY public.category_attribute_mapping
    ADD CONSTRAINT category_attribute_mapping_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.category_attributes(id) ON DELETE CASCADE;

CREATE TRIGGER tr_update_category_attribute_sort_order BEFORE INSERT ON public.category_attribute_mapping FOR EACH ROW EXECUTE FUNCTION public.update_category_attribute_sort_order();