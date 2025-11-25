<?php
function handleSubCategory($method, $body) {
    switch ($method) {
        case 'GET':
            if (isset($_GET['id'])) {
                $query = "SELECT sc.id, sc.name, sc.category_id, c.name as category_name 
                          FROM sub_category sc 
                          JOIN category c ON sc.category_id = c.id 
                          WHERE sc.id = ? AND sc.deleted_at IS NULL";
                $data = getData($query, 'i', [$_GET['id']]);
                
                if (empty($data)) {
                    echo json_encode(['message' => 'Alkategória nem található!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(404);
                }
                echo json_encode($data[0], JSON_UNESCAPED_UNICODE);
            } else {
                $query = "SELECT sc.id, sc.name, sc.category_id, c.name as category_name 
                          FROM sub_category sc 
                          JOIN category c ON sc.category_id = c.id 
                          WHERE sc.deleted_at IS NULL 
                          ORDER BY sc.name";
                $data = getData($query);
                echo json_encode($data, JSON_UNESCAPED_UNICODE);
            }
            break;
        case 'POST':
            if (empty($body['name']) || empty($body['category_id'])) {
                echo json_encode(['message' => 'Hiányzó adat: név és kategória kötelező!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $existing = getData("SELECT id FROM sub_category WHERE name = ? AND deleted_at IS NULL", 's', [$body['name']]);
            if (!empty($existing)) {
                echo json_encode(['message' => 'Már létezik ilyen nevű alkategória!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $success = changeData("INSERT INTO sub_category (name, category_id) VALUES (?, ?)", 'si', [$body['name'], $body['category_id']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            echo json_encode(['message' => 'Alkategória létrehozva'], JSON_UNESCAPED_UNICODE);
            http_response_code(201);
            break;
        case 'PUT':
            if (empty($body['id']) || empty($body['name']) || empty($body['category_id'])) {
                echo json_encode(['message' => 'Hiányzó adat!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $existing = getData("SELECT id FROM sub_category WHERE name = ? AND id != ? AND deleted_at IS NULL", 'si', [$body['name'], $body['id']]);
            if (!empty($existing)) {
                echo json_encode(['message' => 'Ez a név már foglalt!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $success = changeData("UPDATE sub_category SET name = ?, category_id = ? WHERE id = ?", 'sii', [$body['name'], $body['category_id'], $body['id']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            echo json_encode(['message' => 'Sikeres frissítés'], JSON_UNESCAPED_UNICODE);
            break;
        case 'DELETE':
            if (empty($_GET['id'])) {
                echo json_encode(['message' => 'Hiányzó ID!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $id = $_GET['id'];
            
            $dependencies = getData("SELECT id FROM product_type WHERE sub_id = ? AND deleted_at IS NULL", 'i', [$id]);
            if (!empty($dependencies)) {
                echo json_encode(['message' => 'Nem törölhető: Ehhez az alkategóriához aktív terméktípusok tartoznak!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $success = changeData("UPDATE sub_category SET deleted_at = CURDATE() WHERE id = ?", 'i', [$id]);
            
            if (!$success) {
                echo json_encode(['message' => 'Adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }
            
            http_response_code(204);
            break;
        default:
            http_response_code(405);
            break;
    }
}