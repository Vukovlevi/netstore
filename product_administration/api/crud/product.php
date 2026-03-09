<?php
function handleProduct($method, $body) {
    header('Content-Type: application/json; charset=utf-8');

    $sanitizeDate = function($dateInput) {
        if (empty($dateInput)) return null;
        $d = trim((string)$dateInput);
        
        if ($d === 'null' || $d === 'undefined') {
            return null;
        }

        if (strlen($d) === 4 && ctype_digit($d)) {
            return $d . '-01-01';
        }

        if (strpos($d, 'T') !== false) {
            $parts = explode('T', $d);
            $d = $parts[0];
        }

        if (preg_match('/^\d{4}-\d{2}-\d{2}$/', $d)) {
            return $d;
        }

        return null;
    };

    try {
        switch ($method) {
            case 'GET':
                if (isset($_GET['id'])) {
                    $query = "SELECT product.id, product.name, product.description, product.amount, product.size, product.size_type, product.expires_at, product.price, product.discount, product.warranty, product.type_id, product.brand_id, 
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
                    $query = "SELECT product.id, product.name, product.description, product.amount, product.size, product.size_type, product.expires_at, product.price, product.discount, product.warranty, product.type_id, product.brand_id, 
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
                if (empty($body['name']) || empty($body['description']) || !isset($body['amount']) || empty($body['size']) || empty($body['size_type']) || !isset($body['price']) || !isset($body['discount']) || empty($body['type_id']) || empty($body['brand_id'])) {
                    echo json_encode(['message' => 'Hiányzó adat: név, leírás, mennyiség, kiszerelés, ár, kedvezmény, típus és márka kötelező!'], JSON_UNESCAPED_UNICODE);
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

                $warrantyDate = isset($body['warranty']) ? $sanitizeDate($body['warranty']) : null;
                $expiresDate = isset($body['expires_at']) ? $sanitizeDate($body['expires_at']) : null;

                $success = changeData("INSERT INTO product (name, description, amount, size, size_type, expires_at, price, discount, warranty, type_id, brand_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", 'ssisssidsii', [
                    $body['name'], 
                    $body['description'],
                    (int)$body['amount'],
                    $body['size'],
                    $body['size_type'],
                    $expiresDate,
                    (int)$body['price'],
                    (int)$body['discount']/100,
                    $warrantyDate,
                    (int)$body['type_id'], 
                    (int)$body['brand_id']
                ]);
                
                if ($success !== true && $success !== 1 && is_string($success)) {
                    echo json_encode(['message' => 'Adatbázis hiba: ' . $success], JSON_UNESCAPED_UNICODE);
                    return http_response_code(500);
                }
                
                if (!$success) {
                    echo json_encode(['message' => 'Ismeretlen hiba a létrehozáskor!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(500);
                }

                $newProduct = getData("SELECT id, name FROM product WHERE name = ? AND brand_id = ? AND deleted_at IS NULL ORDER BY id DESC LIMIT 1", 'si', [$body['name'], $body['brand_id']]);
                
                if (!empty($newProduct)) {
                     echo json_encode($newProduct[0], JSON_UNESCAPED_UNICODE);
                } else {
                     echo json_encode(['message' => 'Termék létrehozva', 'name' => $body['name']], JSON_UNESCAPED_UNICODE);
                }
                http_response_code(201);
                break;

            case 'PUT':
                if (empty($body['id']) || empty($body['name']) || empty($body['description']) || !isset($body['amount']) || empty($body['size']) || empty($body['size_type']) || !isset($body['price']) || !isset($body['discount']) || empty($body['type_id']) || empty($body['brand_id'])) {
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

                $warrantyDate = isset($body['warranty']) ? $sanitizeDate($body['warranty']) : null;
                $expiresDate = isset($body['expires_at']) ? $sanitizeDate($body['expires_at']) : null;

                $success = changeData("UPDATE product SET name = ?, description = ?, amount = ?, size = ?, size_type = ?, expires_at = ?, price = ?, discount = ?, warranty = ?, type_id = ?, brand_id = ? WHERE id = ?", 'ssisssidisii', [
                    $body['name'], 
                    $body['description'],
                    (int)$body['amount'],
                    $body['size'],
                    $body['size_type'],
                    $expiresDate,
                    (int)$body['price'],
                    (int)$body['discount']/100,
                    $warrantyDate,
                    (int)$body['type_id'], 
                    (int)$body['brand_id'], 
                    (int)$body['id']
                ]);
                
                if ($success !== true && $success !== 1 && is_string($success)) {
                    echo json_encode(['message' => 'Adatbázis hiba: ' . $success], JSON_UNESCAPED_UNICODE);
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
                
                if ($success !== true && $success !== 1 && is_string($success)) {
                    echo json_encode(['message' => 'Adatbázis hiba: ' . $success], JSON_UNESCAPED_UNICODE);
                    return http_response_code(500);
                }
                
                http_response_code(204);
                break;

            default:
                return http_response_code(405);
        }
    } catch (Exception $e) {
        http_response_code(500);
        echo json_encode(['message' => 'Szerver hiba: ' . $e->getMessage()], JSON_UNESCAPED_UNICODE);
    }
}