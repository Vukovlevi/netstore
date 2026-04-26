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
                if (isset($_GET['distinct_size_types'])) {
                    $rows = getData("SELECT DISTINCT size_type FROM product WHERE deleted_at IS NULL AND size_type IS NOT NULL AND size_type != '' ORDER BY size_type");
                    $sizeTypes = array_map(function($r) { return $r['size_type']; }, $rows ?: []);
                    echo json_encode($sizeTypes, JSON_UNESCAPED_UNICODE);
                    break;
                }
                if (isset($_GET['check_deleted']) && isset($_GET['brand_id'])) {
                    $needle = trim((string)$_GET['check_deleted']);
                    if ($needle === '') {
                        echo json_encode(null, JSON_UNESCAPED_UNICODE);
                        break;
                    }
                    $row = getData("SELECT id, name, brand_id, deleted_at FROM product WHERE LOWER(TRIM(name)) = LOWER(?) AND brand_id = ? AND deleted_at IS NOT NULL ORDER BY id DESC LIMIT 1", 'si', [$needle, (int)$_GET['brand_id']]);
                    echo json_encode(!empty($row) ? $row[0] : null, JSON_UNESCAPED_UNICODE);
                    break;
                }
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
                        http_response_code(404);
                        echo json_encode(['message' => 'Termék nem található vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                        return;
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
                $nameInput = isset($body['name']) ? trim((string)$body['name']) : '';

                if (!empty($body['restore']) && !empty($body['id']) && $nameInput !== '' && !empty($body['brand_id']) && !empty($body['type_id'])) {
                    $deletedRow = getData("SELECT id FROM product WHERE id = ? AND deleted_at IS NOT NULL", 'i', [$body['id']]);

                    if (empty($deletedRow)) {
                        http_response_code(404);
                        echo json_encode(['message' => 'Nincs visszaállítható termék ezzel az azonosítóval!'], JSON_UNESCAPED_UNICODE);
                        return;
                    }

                    $clash = getData("SELECT id FROM product WHERE LOWER(TRIM(name)) = LOWER(?) AND brand_id = ? AND id != ? AND deleted_at IS NULL", 'sii', [$nameInput, $body['brand_id'], $body['id']]);

                    if (!empty($clash)) {
                        http_response_code(409);
                        echo json_encode(['message' => 'Már létezik ilyen nevű aktív termék ezen a márkán belül!'], JSON_UNESCAPED_UNICODE);
                        return;
                    }

                    $warrantyDate = isset($body['warranty']) ? $sanitizeDate($body['warranty']) : null;
                    $expiresDate = isset($body['expires_at']) ? $sanitizeDate($body['expires_at']) : null;

                    $restored = changeData("UPDATE product SET name = ?, description = ?, amount = ?, size = ?, size_type = ?, expires_at = ?, price = ?, discount = ?, warranty = ?, type_id = ?, brand_id = ?, deleted_at = NULL WHERE id = ?", 'ssisssidssii', [
                        $nameInput,
                        $body['description'],
                        (int)$body['amount'],
                        $body['size'],
                        $body['size_type'],
                        $expiresDate,
                        (int)$body['price'],
                        (float)$body['discount'],
                        $warrantyDate,
                        (int)$body['type_id'],
                        (int)$body['brand_id'],
                        (int)$body['id']
                    ]);

                    if ($restored !== true && $restored !== 1 && is_string($restored)) {
                        http_response_code(500);
                        echo json_encode(['message' => 'Adatbázis hiba: ' . $restored], JSON_UNESCAPED_UNICODE);
                        return;
                    }

                    if (!$restored) {
                        http_response_code(500);
                        echo json_encode(['message' => 'Ismeretlen hiba a visszaállításkor!'], JSON_UNESCAPED_UNICODE);
                        return;
                    }

                    $row = getData("SELECT id, name FROM product WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
                    http_response_code(200);
                    if (!empty($row)) {
                        echo json_encode($row[0], JSON_UNESCAPED_UNICODE);
                    } else {
                        echo json_encode(['message' => 'Termék visszaállítva', 'name' => $nameInput], JSON_UNESCAPED_UNICODE);
                    }
                    break;
                }

                if ($nameInput === '' || empty($body['description']) || !isset($body['amount']) || empty($body['size']) || empty($body['size_type']) || !isset($body['price']) || !isset($body['discount']) || empty($body['type_id']) || empty($body['brand_id'])) {
                    http_response_code(400);
                    echo json_encode(['message' => 'Hiányzó adat: név, leírás, mennyiség, kiszerelés, ár, kedvezmény, típus és márka kötelező!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $existing = getData("SELECT id, deleted_at FROM product WHERE LOWER(TRIM(name)) = LOWER(?) AND brand_id = ? ORDER BY (deleted_at IS NULL) DESC LIMIT 1", 'si', [$nameInput, $body['brand_id']]);

                if (!empty($existing)) {
                    if (is_null($existing[0]['deleted_at'])) {
                        http_response_code(409);
                        echo json_encode(['message' => 'Már létezik ilyen nevű termék ezen a márkán belül!'], JSON_UNESCAPED_UNICODE);
                        return;
                    } else {
                        http_response_code(409);
                        echo json_encode([
                            'message' => 'Létezik egy korábban törölt termék ezen a néven és márkán belül. Szeretné visszaállítani?',
                            'restorable' => true,
                            'id' => (int)$existing[0]['id']
                        ], JSON_UNESCAPED_UNICODE);
                        return;
                    }
                }

                $warrantyDate = isset($body['warranty']) ? $sanitizeDate($body['warranty']) : null;
                $expiresDate = isset($body['expires_at']) ? $sanitizeDate($body['expires_at']) : null;

                $success = changeData("INSERT INTO product (name, description, amount, size, size_type, expires_at, price, discount, warranty, type_id, brand_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", 'ssisssidsii', [
                    $nameInput,
                    $body['description'],
                    (int)$body['amount'],
                    $body['size'],
                    $body['size_type'],
                    $expiresDate,
                    (int)$body['price'],
                    (float)$body['discount'],
                    $warrantyDate,
                    (int)$body['type_id'],
                    (int)$body['brand_id']
                ]);

                if ($success !== true && $success !== 1 && is_string($success)) {
                    http_response_code(500);
                    echo json_encode(['message' => 'Adatbázis hiba: ' . $success], JSON_UNESCAPED_UNICODE);
                    return;
                }

                if (!$success) {
                    http_response_code(500);
                    echo json_encode(['message' => 'Ismeretlen hiba a létrehozáskor!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $newProduct = getData("SELECT id, name FROM product WHERE name = ? AND brand_id = ? AND deleted_at IS NULL ORDER BY id DESC LIMIT 1", 'si', [$nameInput, $body['brand_id']]);

                http_response_code(201);
                if (!empty($newProduct)) {
                     echo json_encode($newProduct[0], JSON_UNESCAPED_UNICODE);
                } else {
                     echo json_encode(['message' => 'Termék létrehozva', 'name' => $nameInput], JSON_UNESCAPED_UNICODE);
                }
                break;

            case 'PUT':
                if (empty($body['id']) || empty($body['name']) || empty($body['description']) || !isset($body['amount']) || empty($body['size']) || empty($body['size_type']) || !isset($body['price']) || !isset($body['discount']) || empty($body['type_id']) || empty($body['brand_id'])) {
                    http_response_code(400);
                    echo json_encode(['message' => "Hiányzó adat: név, leírás, mennyiség, kiszerelés, ár, kedvezmény, típus és márka kötelező!"], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $nameInput = trim((string)$body['name']);

                $prodToUpdate = getData("SELECT id FROM product WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);

                if (empty($prodToUpdate)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Nincs ilyen azonosítójú termék, vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $existingName = getData("SELECT id FROM product WHERE LOWER(TRIM(name)) = LOWER(?) AND brand_id = ? AND id != ? AND deleted_at IS NULL", 'sii', [$nameInput, $body['brand_id'], $body['id']]);

                if (!empty($existingName)) {
                    http_response_code(409);
                    echo json_encode(['message' => 'Már létezik ilyen nevű termék ezen a márkán belül!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $warrantyDate = isset($body['warranty']) ? $sanitizeDate($body['warranty']) : null;
                $expiresDate = isset($body['expires_at']) ? $sanitizeDate($body['expires_at']) : null;

                $success = changeData("UPDATE product SET name = ?, description = ?, amount = ?, size = ?, size_type = ?, expires_at = ?, price = ?, discount = ?, warranty = ?, type_id = ?, brand_id = ? WHERE id = ?", 'ssisssidssii', [
                    $nameInput,
                    $body['description'],
                    (int)$body['amount'],
                    $body['size'],
                    $body['size_type'],
                    $expiresDate,
                    (int)$body['price'],
                    (float)$body['discount'],
                    $warrantyDate,
                    (int)$body['type_id'],
                    (int)$body['brand_id'],
                    (int)$body['id']
                ]);

                if ($success !== true && $success !== 1 && is_string($success)) {
                    http_response_code(500);
                    echo json_encode(['message' => 'Adatbázis hiba: ' . $success], JSON_UNESCAPED_UNICODE);
                    return;
                }

                echo json_encode(['message' => 'Sikeres frissítés'], JSON_UNESCAPED_UNICODE);
                break;

            case 'DELETE':
                if (empty($_GET['id'])) {
                    http_response_code(400);
                    echo json_encode(['message' => "Hiányzó adat: az azonosító megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $prodToDelete = getData("SELECT id FROM product WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

                if (empty($prodToDelete)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Nincs ilyen azonosítójú termék!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $success = changeData("UPDATE product SET deleted_at = CURDATE() WHERE id = ?", 'i', [$_GET['id']]);

                if ($success !== true && $success !== 1 && is_string($success)) {
                    http_response_code(500);
                    echo json_encode(['message' => 'Adatbázis hiba: ' . $success], JSON_UNESCAPED_UNICODE);
                    return;
                }

                http_response_code(204);
                break;

            default:
                http_response_code(405);
                return;
        }
    } catch (Exception $e) {
        http_response_code(500);
        echo json_encode(['message' => 'Szerver hiba: ' . $e->getMessage()], JSON_UNESCAPED_UNICODE);
    }
}
