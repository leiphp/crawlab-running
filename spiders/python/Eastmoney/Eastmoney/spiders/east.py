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
        pass
