<?php
function handleSearchProduct($method, $body) {
    header('Content-Type: application/json; charset=utf-8');

    if ($method === 'POST') {
        try {
            $page = isset($body['page']) ? (int)$body['page'] : 1;
            if ($page < 1) $page = 1;
            $limit = 25;
            $offset = ($page - 1) * $limit;

            $baseSql = "FROM product p
                        LEFT JOIN brand b ON p.brand_id = b.id
                        LEFT JOIN product_type pt ON p.type_id = pt.id
                        LEFT JOIN sub_category sc ON pt.sub_id = sc.id
                        LEFT JOIN category c ON sc.category_id = c.id
                        LEFT JOIN storing_condition st ON pt.storing_condition_id = st.id
                        WHERE p.deleted_at IS NULL";

            $types = "";
            $params = [];

            if (!empty($body['name'])) {
                $baseSql .= " AND p.name LIKE ?";
                $types .= "s";
                $params[] = "%" . $body['name'] . "%";
            }

            if (!empty($body['description'])) {
                $baseSql .= " AND p.description LIKE ?";
                $types .= "s";
                $params[] = "%" . $body['description'] . "%";
            }

            if (!empty($body['category_id'])) {
                $baseSql .= " AND c.id = ?";
                $types .= "i";
                $params[] = (int)$body['category_id'];
            }

            if (!empty($body['sub_category_id'])) {
                $baseSql .= " AND sc.id = ?";
                $types .= "i";
                $params[] = (int)$body['sub_category_id'];
            }

            if (!empty($body['type_id'])) {
                $baseSql .= " AND pt.id = ?";
                $types .= "i";
                $params[] = (int)$body['type_id'];
            }

            if (!empty($body['brand_id'])) {
                $baseSql .= " AND b.id = ?";
                $types .= "i";
                $params[] = (int)$body['brand_id'];
            }

            if (!empty($body['storing_condition_id'])) {
                $baseSql .= " AND st.id = ?";
                $types .= "i";
                $params[] = (int)$body['storing_condition_id'];
            }

            if (isset($body['amount_min']) && $body['amount_min'] !== '') {
                $baseSql .= " AND p.amount >= ?";
                $types .= "i";
                $params[] = (int)$body['amount_min'];
            }
            if (isset($body['amount_max']) && $body['amount_max'] !== '') {
                $baseSql .= " AND p.amount <= ?";
                $types .= "i";
                $params[] = (int)$body['amount_max'];
            }

            if (isset($body['price_min']) && $body['price_min'] !== '') {
                $baseSql .= " AND p.price >= ?";
                $types .= "i";
                $params[] = (int)$body['price_min'];
            }
            if (isset($body['price_max']) && $body['price_max'] !== '') {
                $baseSql .= " AND p.price <= ?";
                $types .= "i";
                $params[] = (int)$body['price_max'];
            }

            if (isset($body['size_val']) && $body['size_val'] !== '') {
                $baseSql .= " AND p.size = ?";
                $types .= "d";
                $params[] = (float)$body['size_val'];
            }

            if (!empty($body['size_type'])) {
                $baseSql .= " AND p.size_type = ?";
                $types .= "s";
                $params[] = $body['size_type'];
            }

            if (isset($body['is_discounted']) && $body['is_discounted'] === 'true') {
                $baseSql .= " AND p.discount > 0";
            }

            if (isset($body['has_warranty']) && $body['has_warranty'] === 'true') {
                $baseSql .= " AND p.warranty IS NOT NULL AND p.warranty > CURDATE()";
            }

            if (isset($body['show_expired']) && $body['show_expired'] === 'true') {
                $baseSql .= " AND p.expires_at IS NOT NULL AND p.expires_at < CURDATE()";
            }

            $countSql = "SELECT COUNT(*) as total " . $baseSql;
            if (strlen($types) > 0) {
                $countResult = getData($countSql, $types, $params);
            } else {
                $countResult = getData($countSql);
            }
            $total = isset($countResult[0]['total']) ? (int)$countResult[0]['total'] : 0;

            $sql = "SELECT p.*,
                           b.name as brand_name,
                           pt.name as type_name,
                           sc.name as sub_category_name,
                           c.name as category_name,
                           st.description as storing_condition_name
                    " . $baseSql . " ORDER BY p.name ASC LIMIT ? OFFSET ?";

            $types .= "ii";
            $params[] = $limit;
            $params[] = $offset;

            $results = getData($sql, $types, $params);
            if (!is_array($results)) {
                $results = [];
            }

            echo json_encode([
                'data' => $results,
                'total' => $total,
                'page' => $page,
                'limit' => $limit
            ], JSON_UNESCAPED_UNICODE);

        } catch (Exception $e) {
            http_response_code(500);
            echo json_encode(['message' => 'Szerver hiba: ' . $e->getMessage()], JSON_UNESCAPED_UNICODE);
        }
    } else {
        http_response_code(405);
    }
}