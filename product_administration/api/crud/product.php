<?php
function handleProduct($method, $body) {
    switch ($method) {
        case 'GET':
            if (isset($_GET['id'])) {
                $query = "SELECT product.id, product.name, product.type_id, product.brand_id, 
                                 product_type.name as type_name, brand.name as brand_name, 
                                 product_type.sub_id, sub_category.category_id
                          FROM product 
                          LEFT JOIN product_type ON product.type_id = product_type.id 
                          LEFT JOIN brand ON product.brand_id = brand.id 
                          LEFT JOIN sub_category ON product_type.sub_id = sub_category.id
                          WHERE product.id = ? AND product.deleted_at IS NULL";
                $product = getData($query, 'i', [$_GET['id']]);
                
                if (empty($product)) {
                    echo json_encode(['message' => 'Termék nem található vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(404);
                }
                echo json_encode($product[0], JSON_UNESCAPED_UNICODE);
            } else {
                $query = "SELECT product.id, product.name, product.type_id, product.brand_id, 
                                 product_type.name as type_name, brand.name as brand_name 
                          FROM product 
                          LEFT JOIN product_type ON product.type_id = product_type.id 
                          LEFT JOIN brand ON product.brand_id = brand.id 
                          WHERE product.deleted_at IS NULL 
                          ORDER BY product.name";
                $products = getData($query);
                echo json_encode($products, JSON_UNESCAPED_UNICODE);
            }
            break;
        case 'POST':
            if (empty($body['name']) || empty($body['type_id']) || empty($body['brand_id'])) {
                echo json_encode(['message' => 'Hiányzó adat: név, típus és márka kötelező!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $existing = getData("SELECT id, deleted_at FROM product WHERE name = ? AND brand_id = ?", 'si', [$body['name'], $body['brand_id']]);
            
            if (!empty($existing)) {
                if (is_null($existing[0]['deleted_at'])) {
                    echo json_encode(['message' => 'Már létezik ilyen nevű aktív termék ezen a márkán belül!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                } else {
                    echo json_encode(['message' => 'Már létezik ilyen termék, de törölve lett.'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                }
            }

            $success = changeData("INSERT INTO product (name, type_id, brand_id) VALUES (?, ?, ?)", 'sii', [$body['name'], $body['type_id'], $body['brand_id']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            $newProduct = getData("SELECT id, name FROM product WHERE name = ? AND brand_id = ? AND deleted_at IS NULL ORDER BY id DESC LIMIT 1", 'si', [$body['name'], $body['brand_id']]);
            
            if (empty($newProduct)) {
                 echo json_encode(['message' => 'Termék létrehozva', 'name' => $body['name']], JSON_UNESCAPED_UNICODE);
            } else {
                 echo json_encode($newProduct[0], JSON_UNESCAPED_UNICODE);
            }
            http_response_code(201);
            break;
        case 'PUT':
            if (empty($body['id']) || empty($body['name']) || empty($body['type_id']) || empty($body['brand_id'])) {
                echo json_encode(['message' => "Hiányzó adat!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $prodToUpdate = getData("SELECT id FROM product WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            
            if (empty($prodToUpdate)) {
                echo json_encode(['message' => 'Nincs ilyen azonosítójú termék, vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }

            $existingName = getData("SELECT id FROM product WHERE name = ? AND brand_id = ? AND id != ?", 'sii', [$body['name'], $body['brand_id'], $body['id']]);
            
            if (!empty($existingName)) {
                echo json_encode(['message' => 'Ez a név már foglalt ennél a márkánál!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $success = changeData("UPDATE product SET name = ?, type_id = ?, brand_id = ? WHERE id = ?", 'siii', [$body['name'], $body['type_id'], $body['brand_id'], $body['id']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            echo json_encode(['message' => 'Sikeres frissítés'], JSON_UNESCAPED_UNICODE);
            break;
        case 'DELETE':
            if (empty($_GET['id'])) {
                echo json_encode(['message' => "Hiányzó adat: 'id'!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $prodToDelete = getData("SELECT id FROM product WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
            
            if (empty($prodToDelete)) {
                echo json_encode(['message' => 'Nincs ilyen azonosítójú termék!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }

            $success = changeData("UPDATE product SET deleted_at = CURDATE() WHERE id = ?", 'i', [$_GET['id']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }
            
            http_response_code(204);
            break;
        default:
            return http_response_code(405);
    }
}