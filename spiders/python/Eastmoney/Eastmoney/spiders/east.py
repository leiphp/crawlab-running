# -*- coding: utf-8 -*-
import scrapy
from ..items import EastmoneyItem


class EastSpider(scrapy.Spider):
    name = 'east'
    allowed_domains = ['finance.eastmoney.com']
    start_urls = ['https://finance.eastmoney.com/']

    def parse(self, response):
        items=EastmoneyItem()
        lists=response.xpath('//div[@class="left"]')
        for i in lists:
            items['name']=i.xpath('./a/@title').get()
            items['url']=i.xpath('./a//@href').get()

            yield items
#             yield scrapy.Request(items['url'],callback=self.parse_detail,meta={"items": items})


    # 解析详情页
#     def parse_detail(self, response):
#         items = response.meta["items"]
#         # 获取详情页的内容、图片
#         items["content"] = response.xpath('//*[@id="ContentBody"]').extract()
#         yield items  # 对返回的数据进行处理
