-- Views migration

CREATE VIEW public.storefront_orders_view AS
 SELECT storefront_orders.id,
    storefront_orders.order_number,
    storefront_orders.storefront_id,
    storefront_orders.customer_id,
    storefront_orders.payment_transaction_id,
    storefront_orders.subtotal_amount AS subtotal,
    storefront_orders.shipping_amount AS shipping,
    storefront_orders.tax_amount AS tax,
    storefront_orders.discount,
    storefront_orders.total_amount AS total,
    storefront_orders.commission_amount,
    storefront_orders.seller_amount,
    storefront_orders.currency,
    storefront_orders.status,
    storefront_orders.escrow_release_date,
    storefront_orders.escrow_days,
    storefront_orders.shipping_address,
    storefront_orders.billing_address,
    storefront_orders.shipping_method,
    storefront_orders.shipping_provider,
    storefront_orders.tracking_number,
    storefront_orders.customer_notes,
    storefront_orders.seller_notes,
    storefront_orders.payment_method,
    storefront_orders.payment_status,
    storefront_orders.notes,
    storefront_orders.metadata,
    storefront_orders.confirmed_at,
    storefront_orders.shipped_at,
    storefront_orders.delivered_at,
    storefront_orders.cancelled_at,
    storefront_orders.created_at,
    storefront_orders.updated_at
   FROM public.storefront_orders;

CREATE VIEW public.v_attribute_groups_with_items AS
 SELECT ag.id AS group_id,
    ag.name AS group_name,
    ag.display_name AS group_display_name,
    ag.icon AS group_icon,
    ag.sort_order AS group_sort_order,
    agi.id AS item_id,
    agi.attribute_id,
    ca.name AS attribute_name,
    ca.display_name AS attribute_display_name,
    agi.icon AS attribute_icon,
    agi.custom_display_name,
    agi.sort_order AS attribute_sort_order
   FROM ((public.attribute_groups ag
     LEFT JOIN public.attribute_group_items agi ON ((ag.id = agi.group_id)))
     LEFT JOIN public.category_attributes ca ON ((agi.attribute_id = ca.id)))
  WHERE (ag.is_active = true)
  ORDER BY ag.sort_order, agi.sort_order;

CREATE VIEW public.v_category_attributes AS
 SELECT cam.category_id,
    cam.attribute_id,
    cam.is_enabled,
    cam.is_required,
    cam.sort_order,
    cam.custom_component,
    ca.name,
    ca.display_name,
    ca.attribute_type,
    ca.options,
    ca.validation_rules,
    ca.is_searchable,
    ca.is_filterable,
    ca.custom_component AS default_custom_component,
    mc.name AS category_name,
    mc.slug AS category_slug
   FROM ((public.category_attribute_mapping cam
     JOIN public.category_attributes ca ON ((cam.attribute_id = ca.id)))
     JOIN public.marketplace_categories mc ON ((cam.category_id = mc.id)))
  ORDER BY cam.category_id, cam.sort_order, ca.sort_order;