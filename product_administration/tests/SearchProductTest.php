<?php

final class SearchProductTest extends HandlerTestCase
{
    public function testNonPostMethodReturns405(): void
    {
        $res = $this->runAndCapture(fn () => handleSearchProduct('GET', null));

        $this->assertSame(405, $res['code']);
    }

    public function testEmptyBodyProducesUnfilteredCountAndQuery(): void
    {
        DbStub::queueGetData(
            [['total' => 0]],
            []
        );

        $res = $this->runAndCapture(fn () => handleSearchProduct('POST', []));

        $this->assertSame(200, $res['code']);
        $this->assertSame(['data' => [], 'total' => 0, 'page' => 1, 'limit' => 25], $res['json']);

        [$countCall, $dataCall] = DbStub::$getDataCalls;

        $this->assertStringStartsWith('SELECT COUNT(*)', $countCall['sql']);
        $this->assertNull($countCall['types'], 'No bind types should be set when no filters exist');

        $this->assertSame('ii', $dataCall['types']);
        $this->assertSame([25, 0], $dataCall['data']);
    }

    public function testNameFilterUsesLikeWithWildcards(): void
    {
        DbStub::queueGetData([['total' => 1]], [['id' => 1]]);

        $this->runAndCapture(fn () => handleSearchProduct('POST', ['name' => 'tej']));

        $count = DbStub::$getDataCalls[0];
        $this->assertStringContainsString('p.name LIKE ?', $count['sql']);
        $this->assertSame('s', $count['types']);
        $this->assertSame(['%tej%'], $count['data']);

        $data = DbStub::$getDataCalls[1];
        $this->assertSame('sii', $data['types']);
        $this->assertSame(['%tej%', 25, 0], $data['data']);
    }

    public function testMultipleFiltersComposeTypesInOrder(): void
    {
        DbStub::queueGetData([['total' => 4]], []);

        $body = [
            'name'         => 'sajt',
            'category_id'  => '5',
            'price_min'    => '100',
            'price_max'    => '500',
            'is_discounted'=> 'true',
        ];
        $this->runAndCapture(fn () => handleSearchProduct('POST', $body));

        $count = DbStub::$getDataCalls[0];
        $this->assertSame('siii', $count['types']);
        $this->assertSame(['%sajt%', 5, 100, 500], $count['data']);
        $this->assertStringContainsString('p.discount > 0', $count['sql']);
    }

    public function testInvalidPageDefaultsToOne(): void
    {
        DbStub::queueGetData([['total' => 0]], []);

        $res = $this->runAndCapture(fn () => handleSearchProduct('POST', ['page' => -3]));

        $this->assertSame(1, $res['json']['page']);
        $dataCall = DbStub::$getDataCalls[1];
        $this->assertSame([25, 0], $dataCall['data']);
    }

    public function testPaginationComputesCorrectOffset(): void
    {
        DbStub::queueGetData([['total' => 200]], []);

        $this->runAndCapture(fn () => handleSearchProduct('POST', ['page' => 4]));

        $dataCall = DbStub::$getDataCalls[1];
        $this->assertSame([25, 75], $dataCall['data'], 'page 4 -> offset 75');
    }


    public function testBooleanFiltersRequireStringTrue(): void
    {
        DbStub::queueGetData([['total' => 0]], []);

        $this->runAndCapture(fn () => handleSearchProduct('POST', ['is_discounted' => true]));
        $this->assertStringNotContainsString('p.discount > 0', DbStub::$getDataCalls[0]['sql']);

        DbStub::reset();
        DbStub::queueGetData([['total' => 0]], []);

        $this->runAndCapture(fn () => handleSearchProduct('POST', ['is_discounted' => 'true']));
        $this->assertStringContainsString('p.discount > 0', DbStub::$getDataCalls[0]['sql']);
    }

    public function testReturnsResultsWithTotal(): void
    {
        DbStub::queueGetData(
            [['total' => 2]],
            [
                ['id' => 1, 'name' => 'A', 'brand_name' => 'Bx'],
                ['id' => 2, 'name' => 'B', 'brand_name' => 'Bx'],
            ]
        );

        $res = $this->runAndCapture(fn () => handleSearchProduct('POST', []));

        $this->assertSame(2, $res['json']['total']);
        $this->assertCount(2, $res['json']['data']);
        $this->assertSame('A', $res['json']['data'][0]['name']);
    }

    public function testSizeValIsBoundAsFloat(): void
    {
        DbStub::queueGetData([['total' => 0]], []);

        $this->runAndCapture(fn () => handleSearchProduct('POST', ['size_val' => '1.5']));

        $count = DbStub::$getDataCalls[0];
        $this->assertSame('d', $count['types']);
        $this->assertSame([1.5], $count['data']);
    }
}
