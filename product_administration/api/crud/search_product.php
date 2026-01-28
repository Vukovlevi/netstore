<?php
function handleSearchProduct($method, $body) {
    header('Content-Type: application/json; charset=utf-8');

    if ($method === 'GET') {
        try {
            $page = isset($_GET['page']) ? (int)$_GET['page'] : 1;
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

            if (!empty($_GET['name'])) {
                $baseSql .= " AND p.name LIKE ?";
                $types .= "s";
                $params[] = "%" . $_GET['name'] . "%";
            }

            if (!empty($_GET['category_id'])) {
                $baseSql .= " AND c.id = ?";
                $types .= "i";
                $params[] = (int)$_GET['category_id'];
            }

            if (!empty($_GET['sub_category_id'])) {
                $baseSql .= " AND sc.id = ?";
                $types .= "i";
                $params[] = (int)$_GET['sub_category_id'];
            }

            if (!empty($_GET['type_id'])) {
                $baseSql .= " AND pt.id = ?";
                $types .= "i";
                $params[] = (int)$_GET['type_id'];
            }

            if (!empty($_GET['brand_id'])) {
                $baseSql .= " AND b.id = ?";
                $types .= "i";
                $params[] = (int)$_GET['brand_id'];
            }

            if (!empty($_GET['storing_condition_id'])) {
                $baseSql .= " AND st.id = ?";
                $types .= "i";
                $params[] = (int)$_GET['storing_condition_id'];
            }

            if (isset($_GET['amount_min']) && $_GET['amount_min'] !== '') {
                $baseSql .= " AND p.amount >= ?";
                $types .= "i";
                $params[] = (int)$_GET['amount_min'];
            }
            if (isset($_GET['amount_max']) && $_GET['amount_max'] !== '') {
                $baseSql .= " AND p.amount <= ?";
                $types .= "i";
                $params[] = (int)$_GET['amount_max'];
            }

            if (isset($_GET['price_min']) && $_GET['price_min'] !== '') {
                $baseSql .= " AND p.price >= ?";
                $types .= "i";
                $params[] = (int)$_GET['price_min'];
            }
            if (isset($_GET['price_max']) && $_GET['price_max'] !== '') {
                $baseSql .= " AND p.price <= ?";
                $types .= "i";
                $params[] = (int)$_GET['price_max'];
            }

            if (isset($_GET['size_val']) && $_GET['size_val'] !== '') {
                $baseSql .= " AND p.size = ?";
                $types .= "d";
                $params[] = (float)$_GET['size_val'];
            }

            if (!empty($_GET['size_type'])) {
                $baseSql .= " AND p.size_type = ?";
                $types .= "s";
                $params[] = $_GET['size_type'];
            }

            if (isset($_GET['is_discounted']) && $_GET['is_discounted'] === 'true') {
                $baseSql .= " AND p.discount > 0";
            }

            if (isset($_GET['has_warranty']) && $_GET['has_warranty'] === 'true') {
                $baseSql .= " AND p.warranty IS NOT NULL AND p.warranty > CURDATE()";
            }

            if (isset($_GET['show_expired']) && $_GET['show_expired'] === 'true') {
                $baseSql .= " AND p.expires_at IS NOT NULL AND p.expires_at < CURDATE()";
            } else {
                $baseSql .= " AND (p.expires_at IS NULL OR p.expires_at >= CURDATE())";
            }

            if (!empty($_GET['other_properties'])) {
                $baseSql .= " AND (p.description LIKE ? OR p.other_properties LIKE ?)";
                $types .= "ss";
                $params[] = "%" . $_GET['other_properties'] . "%";
                $params[] = "%" . $_GET['other_properties'] . "%";
            }

            $countSql = "SELECT COUNT(*) as total " . $baseSql;
            $countResult = getData($countSql, $types, $params);
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