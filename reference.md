# HeroSMS API 参考文档

> 数据来源：HeroSMS API 实时接口，更新时间：2026-05-28

## 获取最新列表

获取完整国家列表：
```
GET https://hero-sms.com/stubs/handler_api.php?action=getCountries&api_key=YOUR_KEY
```

获取完整服务列表：
```
GET https://hero-sms.com/stubs/handler_api.php?action=getServicesList&api_key=YOUR_KEY&lang=cn
```

---

## 国家 ID 对照表

共 194 个可用国家

| ID | 国家 | ID | 国家 | ID | 国家 | ID | 国家 |
|----|------|----|------|----|------|----|------|
| 1 | 乌克兰 | 51 | 白俄罗斯 | 100 | 科威特 | 149 | 索馬里 |
| 2 | 哈萨克斯坦 | 52 | 泰国 | 101 | 薩爾瓦多 | 150 | 剛果 |
| 3 | 中国 | 53 | 沙特阿拉伯 | 102 | 利比亞 | 151 | 智利 |
| 4 | 菲律宾 | 54 | 墨西哥 | 103 | 牙買加 | 152 | 布基納法索 |
| 5 | 缅甸 | 55 | 台湾 | 104 | 特立尼達和多巴哥 | 153 | 黎巴嫩 |
| 6 | 印度尼西亚 | 56 | 西班牙 | 105 | 厄瓜多爾 | 154 | 加蓬 |
| 7 | 马来西亚 | 57 | 伊朗 | 106 | 斯威士蘭 | 155 | 阿爾巴尼亞 |
| 8 | 肯尼亚 | 58 | 阿尔及利亚 | 107 | 阿曼 | 156 | 烏拉圭 |
| 9 | 坦桑尼亚 | 59 | 斯洛文尼亚 | 108 | 波斯尼亞和黑塞哥維那 | 157 | 毛里求斯 |
| 10 | 越南 | 60 | 孟加拉国 | 109 | 多明尼加共和國 | 158 | 丁烷 |
| 11 | 吉尔吉斯斯坦 | 61 | 塞内加尔 | 110 | 敘利亞 | 159 | 马尔代夫 |
| 13 | 以色列 | 62 | 土耳其 | 111 | 卡塔爾 | 160 | 瓜德罗普岛 |
| 14 | 香港 | 63 | 捷克共和国 | 112 | 巴拿馬 | 161 | 土库曼斯坦 |
| 15 | 波兰 | 64 | 斯里兰卡 | 113 | 古巴 | 162 | 法属圭亚那 |
| 16 | 英格兰 | 65 | 秘鲁 | 114 | 毛里塔尼亞 | 163 | 芬兰 |
| 17 | 马达加斯加 | 66 | 巴基斯坦 | 115 | 塞拉利昂 | 164 | 圣卢西亚 |
| 18 | 刚果 | 67 | 新西兰 | 116 | 約旦 | 165 | 卢森堡 |
| 19 | 尼日利亚 | 68 | 几内亚 | 117 | 葡萄牙 | 166 | 圣文森特和格林纳丁斯 |
| 20 | 澳门 | 69 | 马里 | 118 | 巴巴多斯 | 167 | 赤道几内亚 |
| 21 | 埃及 | 70 | 委内瑞拉 | 119 | 布隆迪 | 168 | 吉布地 |
| 22 | 印度 | 71 | 埃塞俄比亚 | 120 | 貝寧 | 169 | 安提瓜和巴布达 |
| 23 | 爱尔兰 | 72 | 蒙古 | 121 | 文萊 | 170 | 开曼群岛 |
| 24 | 柬埔寨 | 73 | 巴西 | 122 | 巴哈馬 | 171 | 黑山共和国 |
| 25 | 老挝 | 74 | 阿富汗 | 123 | 博茨瓦納 | 172 | 丹麥 |
| 26 | 海地 | 75 | 乌干达 | 124 | 伯利茲 | 173 | 瑞士 |
| 27 | 象牙海岸 | 76 | 安哥拉 | 125 | 中非共和國 | 174 | 挪威 |
| 28 | 冈比亚 | 77 | 塞浦路斯 | 126 | 多米尼加 | 175 | 澳大利亞 |
| 29 | 塞尔维亚 | 78 | 法國 | 127 | 格林納達 | 176 | 厄立特里亞 |
| 30 | 也门 | 79 | 巴布亞新幾內亞 | 128 | 佐治亞州 | 177 | 南蘇丹 |
| 31 | 南非 | 80 | 莫桑比克 | 129 | 希臘 | 178 | 聖多美和普林西比 |
| 32 | 罗马尼亚 | 81 | 尼泊爾 | 130 | 幾內亞比紹 | 179 | 阿魯巴島 |
| 33 | 哥伦比亚 | 82 | 比利時 | 131 | 圭亞那 | 180 | 蒙特塞拉特 |
| 34 | 爱沙尼亚 | 83 | 保加利亞 | 132 | 冰島 | 181 | 安圭拉島 |
| 35 | 阿塞拜疆 | 84 | 匈牙利 | 133 | 科摩羅 | 182 | 日本 |
| 36 | 加拿大 | 85 | 摩爾多瓦 | 134 | 聖基茨和尼維斯 | 183 | 北馬其頓 |
| 37 | 摩洛哥 | 86 | 義大利 | 135 | 利比里亞 | 184 | 塞舌爾共和國 |
| 38 | 加纳 | 87 | 巴拉圭 | 136 | 萊索托 | 185 | 新喀裡多尼亞 |
| 39 | 阿根廷 | 88 | 洪都拉斯 | 137 | 馬拉維 | 186 | 佛得角 |
| 40 | 乌兹别克斯坦 | 89 | 突尼斯 | 138 | 納米比亞 | 187 | 美国（物理) |
| 41 | 喀麦隆 | 90 | 尼加拉瓜 | 139 | 尼日爾 | 188 | 巴勒斯坦 |
| 42 | 乍得 | 91 | 東帝汶 | 140 | 盧旺達 | 189 | 斐濟 |
| 43 | 德国 | 92 | 玻利維亞 | 141 | 斯洛伐克 | 196 | 新加坡共和国 |
| 44 | 立陶宛 | 93 | 哥斯達黎加 | 142 | 蘇里南 | 198 | 萨摩亚 |
| 45 | 克罗地亚 | 94 | 危地馬拉 | 143 | 塔吉克斯坦 | 199 | 马耳他 |
| 46 | 瑞典 | 95 | 阿拉伯聯合酋長國 | 144 | 摩納哥 | 201 | 直布罗陀 |
| 47 | 伊拉克 | 96 | 津巴布韋 | 145 | 巴林 | 203 | 科索沃 |
| 48 | 荷兰 | 97 | 波多黎各 | 146 | 團圓 | 204 | 纽埃 |
| 49 | 拉脱维亚 | 98 | 蘇丹蘇丹 | 147 | 贊比亞 | | |
| 50 | 奥地利 | 99 | 多哥 | 148 | 亞美尼亞 | | |

---

## 服务代码对照表

共 805 个服务

| 代码 | 服务名称 | 代码 | 服务名称 |
|------|----------|------|----------|
| full | Full rent | ayb | AsianDating |
| am | Amazon | ayc | HungerStation |
| fb | facebook | ayn | Lime |
| dr | OpenAI | ayo | 360Kredi |
| wa | Whatsapp | ayz | Kimi |
| gp | Ticketmaster | baq | Redbubble |
| go | Google,youtube,Gmail | bav | LeoList |
| wb | WeChat | bay | Genome |
| ig | Instagram+Threads | bba | SwitchUp |
| tg | Telegram | bbu | InternationalCupid |
| ni | Gojek | bcg | 2ГИС |
| tw | Twitter/X | bck | Onion Academy 洋葱学园 |
| ot | Any other | bd | X5ID |
| ka | Shopee | bdd | BetMen |
| ccu | Google Chat | bdt | HOT51 |
| ds | Discord | bdw | Valora  |
| hw | Alipay/Alibaba/1688 | ber | Gumtree |
| cq | Mercado | bfr | Dott |
| ew | Nike | bgj | MoonPay |
| yw | Grindr | bih | Indosaku |
| lf | TikTok/Douyin | bij | Prodege |
| nv | Naver | bix | Hotel101 |
| oi | Tinder | bji | Perfect World 完美世界 |
| mb | Yahoo | bjz | Jeevansathi |
| ub | Uber | bkn | betlive |
| wx | Apple | blr | DocuSign |
| ts | PayPal | blx | 2ememain |
| fr | Dana | blz | MiniPay |
| pm | AOL | bmd | VooV Meeting |
| wr | Walmart | bmi | Sisal |
| im | Imo | bmm | Hey Cash |
| mm | Microsoft | bng | Jush  |
| vi | Viber | bny | Suno |
| qf | RedBook | bod | Genspark |
| fu | Snapchat | boh | Yophone |
| awz | PlayTime | bol | My TELUS |
| abu | BPJSTK | bon | RetailMeNot |
| hb | Twitch | bow | Affirm |
| bcq | Mantan/Kopi Kenangan | box | Thumbtack |
| me | Line messenger | bp | GoFundMe |
| dh | eBay | bpp | Bliq |
| awv | Wallapop | bpq | MChat |
| jg | Grab | bpz | Just Eat |
| ki | 99app | bqb | Weex |
| bw | Signal | bqi | Datanyze |
| acz | Claude  | bqm | Beebs |
| kc | Vinted | brg | Letgo |
| acm | Razer | brr | LemFi |
| ev | Picpay  | bru | myPOS |
| nz | Foodpanda | bry | Mobile DE |
| ang | TOMORO COFFEE | bsm | Ahlan |
| vk | vk.com | bso | PayWell |
| abk | GMX | bss | Trade Republic |
| ju | Indomaret | bsu | Dingtone |
| nf | Netflix | btd | Volcengine 火山引擎 |
| aka | LinkAja | btm | ZeeNow |
| mo | Bumble | btr | Duet |
| alj | Spotify | btx | Amap 高德地图 |
| aez | Shein | bud | WeTalk |
| tn | LinkedIN | bui | OKbet |
| xh | OVO | bvi | Salams |
| aop | Kleinanzeigen | bvs | Вчасно |
| sg | OZON | bwe | Immutable Play |
| bn | Alfagift | bxi | Grupo Madero |
| bzk | Amaze Super App | bxn | Intel SAP |
| ada | TRUTH SOCIAL | bxp | Tennents |
| df | Happn | byd | Airtasker |
| anx | InfinitePay | byp | Kaito |
| xk | DiDi | bzb | Big Cash |
| cn | Fiverr | cac | Vesseo |
| vm | OkCupid | cba | Enilive |
| amb | Vercel | cbh | ZUL |
| sn | OLX | ccl | Cruzeiro |
| nc | Payoneer | cct | Factory |
| dl | Lazada | co | Rediffmail |
| afp | VFS GLOBAL | dy | Zomato |
| pf | pof.com | fs | Şikayet var |
| alg | Ankama | ft | Bookmakers |
| ya | Yandex | gr | Astropay |
| fz | KFC | gs | SamsungShop |
| anh | Cadbury | gt | Gett |
| qj | Whoosh | gx | Hepsiburadacom |
| vz | Hinge | hc | MOMO |
| aaa | Nubank | ib | Immowelt |
| aoy | PLN Mobile | ir | Chispa |
| ac | DoorDash | ji | Monobank |
| pc | Casino/bet/gambling | jx | Swiggy |
| kt | KakaoTalk | kl | kolesa.kz |
| vr | MotorkuX | ls | Careem |
| pd | IFood | ly | Olacabs |
| do | Leboncoin | mx | SoulApp |
| aba | Rappi | ny | BitcoinBon |
| abc | Taptap Send | oc | DealShare |
| afm | myboost | oe | Codashop |
| yy | Venmo | oj | LoveRu |
| bnl | Reddit | pb | SkyTV |
| bz | Blizzard | pu | Justdating |
| ue | Onet | qh | Oriflame |
| als | Greggs  | rc | Skype |
| aik | ZUS Coffee | rs | Lotus |
| xx | Joyride | rt | hily |
| uu | Wildberries | sy | Brahma |
| ok | ok.ru | tc | Rumbler |
| tx | Bolt | ul | Getir |
| uk | Airbnb | uv | BinBin |
| za | JDcom | vd | Betfair |
| hu | Ukrnet | vp | Kwai |
| hx | AliExpress | vs | WinzoGame |
| bxy | MAX | vy | Meta |
| jq | Paysafecard | wc | Craigslist |
| agl | Betano | wd | Столото |
| bgv | Clearpay | ws | Feeld |
| rr | Wolt | yj | eWallet |
| xd | Tokopedia | yk | SportMaster |
| aiw | Roblox | zd | Zilch |
| ann | Bradesco | zl | Airtel |
| gj | Carousell | zr | Papara |
| akp | Her | ve | Dream11 |
| cau | Ero Me | bsy | Air Miles |
| arc | CHECK24 | bvh | ParenTeam |
| bex | Whatnot | alp | Mera Gaon |
| fv | Vidio | boo | Casino Plus |
| bdp | Kredito | jc | IVI |
| kf | Weibo | bnp | Airba Fresh |
| yx | JTExpress | bbd | DBS Bank |
| aff | C6 Bank | bri | Gardrops |
| bfg | EMAG | bhf | InPost |
| ja | Weverse | cr | TenChat |
| btv | Radiate | us | IRCTC |
| zh | Zoho | qr | MEGA |
| abg | PagBank | wu | PrivetMir |
| ais | DiDiFood | bmb | Газпром ID |
| zp | Pinduoduo | blw | BP |
| ajn | Gopuff | tp | IndiaGold |
| auh | KeeTa 美团 | adj | RummyCircle |
| zk | Deliveroo | xy | Depop |
| ah | EscapeFromTarkov | bxw | Credinex |
| azi | TheL | ha | My11Circle |
| bc | GCash | we | DrugVokrug |
| blt | INDOPAKET | bhz | Eurobet |
| ht | Bitso | aym | Langit Musik |
| tm | Akulaku | auo | Mycashbacks |
| apq | WePoker | jy | Sorare |
| ua | BlaBlaCar | ajc | Pochta Rossii |
| abn | Bybit | mi | Zupee |
| bwu | SkyBet | bcy | Finplus |
| zy | Nttgame | bbt | Blink |
| ani | Talabat | md | Banks |
| aje | CupidMedia | agx | MeiQFashion |
| axn | FastMoss | jr | Samokat |
| si | Cita Previa | bkw | Bilkraft |
| aey | Next | ana | Sicredi |
| arf | Enjoei | rn | neftm |
| auc | TotalPass | azz | Crystalbet |
| tl | Truecaller | bly | COSMOTE |
| afz | Klarna | blm | Epic Games  |
| bwv | Manus | hz | Drom |
| lj | Santander | aga | Publi24 |
| adu | Seznam | bte | CROCOBET |
| bbj | OnePay | aih | Fups |
| km | Rozetka | aax | Boyaa |
| afs | Privalia | abl | gpnbonus |
| ahl | Maxim | adc | PlayOJO |
| akr | Voi | ady | TOKYO-CITY |
| aq | Glovo | kw | Foody |
| aqt | Skrill | ms | NovaPoshta |
| ars | Bingo Plus | rj | Detskiy mir |
| bqh | PEDIR GAS | bii | MEXC |
| fk | BLIBLI | atr | RIDE  |
| gf | GoogleVoice | ark | SageMaker Studio Lab |
| mt | Steam | ari | Ring4 |
| bpt | SerpApi  | cbk | Fastwork |
| bwx | Chagee | bbl | Autodesk |
| ex | Linode | bko | Vision Plus |
| kk | Idealista | aqm | Tala |
| qv | Badoo | bsa | will |
| sw | NCsoft | azl | MotoRan  |
| yl | Yalla | ccv | Triumph |
| asb | YUEWEN 阅文集团 | bic | Marktde |
| bbk | FilipinoCupid | bwt | Stanleybet |
| bcx | Bantusaku | azy | Casa Pariurilor |
| wh | TanTan | bou | 汇旺 Huione Pay |
| ane | Supercell | aej | Autoru |
| anm | CaltexGO | brf | Dolap |
| bbq | Chime | agc | VIMpay |
| qq | Tencent QQ | bbf | Neosurf |
| re | Coinbase | bbg | Square |
| sa | AGIBANK | sc | Voggt |
| yi | Yemeksepeti | arr | Linii Lybvi |
| adt | willhaben | bdr | Yara Farmcare |
| ama | WooPlus | tv | Flink |
| axx | Shopback | wp | 163СOM |
| bfw | Finya | bro | YouGov Shopper |
| bl | BIGO LIVE | br | Vkusno i Tochka |
| btn | Itau | bsw | Milanuncios |
| ep | Temu | bqx | Club Q8 |
| oz | Poshmark | bhj | Jòfogàs |
| rl | inDriver | bwo | KabanchikUA |
| zg | Setel | bhl | ati su |
| abo | WEBDE | auj | High 5 Casino |
| afe | GovBr | amy | Otzovik |
| agd | Grailed | btq | ЭЛПЛАТ |
| ahe | Bunq | anw | Premmia |
| alb | Guiche Web | tf | Noon |
| app | ClassPass | apj | Tinkoff VoiceKit |
| bra | Touch n Go TNG | bsl | Oskelly |
| byh | Adobe | aeh | Apteka Aprel |
| tu | Lyft | bls | WeTV |
| bo | Wise | brt | Netwin |
| bxg | BAT | axt | GNJOY |
| cp | Uklon | dt | Delivery Club |
| cw | PaddyPower | boz | Doctu |
| ff | AVON | bxc | DRIVE2 |
| pr | Trendyol | apy | Vinlab |
| sr | Starbucks | yh | hh |
| aaq | Netease | blq | Ari10 |
| aix | Move It | om | Corona |
| ajj | Rebtel | bhb | Wuling |
| akd |  Feels | bhe | Jago |
| aoz | ReclameAQUI | ahk | BlinBeri |
| aqy | SAPO | afr | Ultragaz |
| axj | FIFGROUP MOBILE CUSTOMER | ahr | This Fate |
| brk | Indeed | nt | Sravni |
| dg | Mercari | acb | Spark Driver |
| ij | Revolut | tr | Paysend |
| mw | Transfergo | mg | Magnit |
| te | eFood | blc | Дикси |
| abe | Foodora | axm | Não Me Perturbe |
| agk | Ipsos iSay | brv | 999 md |
| akl | DOKU | ym | youla.ru |
| aok | NETELLER | cj | Dotz |
| aup | Botim | bst | AdmiralBet |
| axr | Match | baj | Bipa |
| bai | SNKRDUNK | of | urent/jet/RuSharing |
| bmj | Betflag | caj | GoMoney |
| dp | ProtonMail | bab | Opera Mini |
| et | Clubhouse | amz | ImmoScout24 |
| li | Baidu | sv | Dostavista |
| ael | Cloud Manager | aaz | Ozan |
| agh | Getnet | avv | shopFarEast |
| ahb | Ubisoft | ky | SpatenOktoberfest |
| ait | FeetFinder | are | Seated |
| akz | Panvel | kv | Rush |
| anl | AttaPoll | bv | Metro |
| aoh | YooMoney | hp | Meesho |
| apd | 2dehands | dj | LUKOIL-AZS |
| avb | Tealive | bwy | Монетка |
| bbm | TrueMoney | zt | Budweiser |
| bhr | Dil Mil | bgt | Alfamidi |
| bli | Scalapay | py | Monese |
| bni | Pets4Homes | apf | Carrefour |
| btl | D4 | ble | Nexon |
| gm | Mocospace | agi | Njuškalo |
| gq | Freelancer | bib | myBCA |
| mv | Fruitz | kq | FotoCasa |
| nq | Trip | ayv | myTVSUPER |
| ns | Oldubil | bog | France Mobilities |
| pz | Lidl | ng | FunPay |
| qx | WorldRemit | bxv | PAPER |
| uf | Eneba | cbs | lalafo |
| ww | BIP | bvu | SwaRail  |
| zo | Kaggle | xu | RecargaPay |
| zs | Bilibili | btu | DDX Fitness |
| aav | Alchemy  | hy | Ininal |
| aex | Neon | aoq | JB Hi-Fi |
| ako | Ryde | jl | Hopi |
| alm | Muzz | ayr | Superbet |
| aof | Guzman y Gomez | afa | CDEK |
| aor | OKX | ahv | Curve |
| arp | Continente | aua | 同程旅行 Tongcheng Travel |
| aub | Smitten | bwl | Guthabende |
| axp | ChargePoint | atu | Sber/Cooper |
| bbr | WalletHub | bxu | Chronodrive |
| bcz | UATAS | akm | LOTTE Mart |
| bfo | MSport | amw | Tise |
| bfv | Chocofamily | fw | 99acres |
| bpc | LOVOO | aab | BharatPe |
| byw | Annoncelight | bke | FOOTDISTRICT |
| bzo | SeLoger | bku | NL International |
| ee | Twilio | bre | Лемана ПРО |
| em | ZéDelivery | bbs | Ecoplatform |
| ie | bet365 | aqq | MPSTATS |
| ip | Burger King | asq | Warpcast |
| ti | cryptocom | bnz | Seven Bet |
| uz | OffGamers | bmx | Canal RCN |
| xr | Tango | acd | Cloud.ru |
| xs | GroupMe | bye | WhatsApp Business |
| xz | paycell | bqr | ASAAS |
| zb | FreeNow | cak | Bazaraki |
| ze | Shpock | acu | CityMall |
| acw | YouDo | bob | Shell GO |
| ahi | Daki | axg | Twin |
| ahj | Strato | bop | Везёт |
| aiq | Prime Opinion | ahn | Riv Gosh |
| aiz | Brevo | acr | QwikCilver |
| alo | Profee | cm | Prom |
| aoe | Sendwave | btp | Связь ON |
| aoi | Cryptonow | aps | Skelbiu |
| aom | Monzo | bxr | AlFursan |
| apl | Sideline | aat | TamTam |
| apr | Capital One | aly | Bebeclub |
| aux | Lightning AI | abb | Coca-Cola |
| avj | SumUp  | bdh | Singledk |
| bbi | aiqfome | bdb | ОБИ |
| bdg | HUD | auf | PREMIER |
| bem | Shafa | bvj | Санофи |
| bgn | Veeka | ae | myGLO |
| bkd | Sahibinden | bsv | AmarthaFin |
| bkv | Zen | cbb | KeepCalling |
| blo | SiliconFlow 硅基流动 | rd | Lenta |
| bme | myIM3 | bdo | AdaModal |
| bnu | Qpon | bje | Meest |
| bos | Casino Portugal | ta | Wink |
| bpe | Сільпо | bwf | MyGo |
| bpi | Askable | ahx | Bitrue |
| bqn | Upward | ccw | 夸克  Quark |
| bqo |  Caffe Nero | bxs | LikeCard |
| bqp | Zara | qg | MoneyPay |
| bsf | Rupiah Cepat | bif | OPAP |
| bwa | Hostinger | ban | BLINK by BonusLink |
| bwd | Fetch | can | BAGI |
| byf | SeaBank | adw | Profi |
| bzv | TOP | cam | Eleme 饿了么 |
| cb | Bazos | bqe | Lottomatica |
| ef | Nextdoor | bnt | Treasury |
| ej | MrQ | buy | GetBlogger |
| fd | Mamba | bas | Stoiximan |
| fh | Lalamove | afu | VseInstrumenty |
| ge | Paytm | ke | Eldorado |
| lc | Subito | aer | PlayerAuctions |
| ma | Mail.ru | aeg | Flowwow |
| my | CAIXA | sl | robota.ua |
| nh | AlloBank | aes | Zolotoye Yabloko |
| qe | GG | zn | Biedronka |
| qo | Moneylion | abh | UOL Host |
| rk | Fotka | sd | dodopizza |
| ry | McDonalds | js | GolosZa |
| vg | ShellBox | aid | Kwiff |
| wg | Skout | acn | Gringo |
| yn | Allegro | arb | Las Vegas Casino |
| yu | Xiaomi | bnf | ContiCazino |
| zm | OfferUp | ahc | START |
| aag | Pockit | kj | YAPPY |
| abd | BeBoo | ov | Beget |
| abq | Upwork | azd | Вплюсе |
| abt | ArenaPlus | bxj | Quero-Quero PAG |
| acc |  LuckyLand Slots | bqy | SNAI |
| ace | Tata Neu | bga | LPTracker |
| aci | Colombian Cupid | buj | ASVLA |
| adi | Zepto | avu | Karos |
| adp | Cabify | adv | Cian |
| aem | AstraPay | bjq | Surveyon  |
| aeq | Godrej | ami | TEAMORU |
| aet | Greywoods | bhp | PPC |
| aeu | TheFork | il | IQOS |
| afc | Bunda | bau | FieldStar |
| afd | Astra Otoshop | bfa | Webmotors  |
| afn | roomster | bvf | GymPlius |
| agb | Smiles | bxf | Disney Plus |
| agg | OneForma | lm | FarPost |
| agj | Marktplaats | arx | Zinia |
| agm | CMB | sh | Vkusvill |
| agv | DoneDeal | blp | MEEFF |
| agw | Adverts | bkk | Poparide |
| ahd | Midnite | bjp | Ria |
| ahf | GAIL`s Bakery | xt | Flipkart |
| aig | Five Surveys | gk | AptekaRU |
| ail | Smotryoshka | pv | Koshelek |
| aim | Firebase | bcr | ESX |
| aiv | Remitly | cy | PSA |
| ajq | MyValue | afk | Chevron |
| aju | Daya Auto | awg | Natura Avon |
| ajv | ShareParty | alx | NutriClub |
| ajy | All Access | ui | RuTube |
| akj | Easycash | bpg | НаПоправку |
| anb | Abastece-aí | st | Auchan |
| anj | Gemini | aos | Loloo |
| ano | Shopify | blh | Winner |
| anq | Hitnspin | afo | KION  |
| aol | Paysera | azn | UFC |
| aon | Binance | bkz | Rewardy |
| apb | eToro | bns | PAC Cash |
| apg | Damai | aow | Geekay |
| api | KKTIX | azs | Магнолия |
| apo | Netmarble | aay | JioMart |
| aqj | BigBasket | xm | Letual |
| aru | UNNI 강남언니 | aoo | Pegasus Airlines |
| asf | TextFree | qy | Zhihu |
| asj | T-online | bit | Домклик |
| asp | PhonePe | bnd | КуулКлевер |
| asy | Fore Coffee | tk | MVideo |
| atl | Watsons MY | sb | Lamoda |
| atn | Fawry | da | MTS CashBack |
| atp | Vonage | aox | Aukro |
| aul | Alignable | brc | Casa it |
| auz | Outlier | bzf | Dabble |
| avk | Quoka  | wt | IZI |
| avt | HuntTheMouse | bla | FDJ ParionsSport |
| awq | Atlas Earth | aug | Magnit Market |
| awr | KCEX | bki | MyACUVUE |
| aws | 7-Eleven | adz | Shokoladnica |
| awu | MosGram | tz | Leyka |
| axs | AIS PLAY |  |  |

---

## API 端点速查

### 查余额
```
GET https://hero-sms.com/stubs/handler_api.php?action=getBalance&api_key=KEY
```
返回：`ACCESS_BALANCE:100.5`

### 报价查询（监控用）
```
GET https://hero-sms.com/stubs/handler_api.php?action=getPrices&api_key=KEY&service=tg&country=6
```
返回：`{country_id: {service_id: {cost, count, physicalCount}}}`

### 获取号码 V2
```
GET https://hero-sms.com/stubs/handler_api.php?action=getNumberV2&api_key=KEY&service=tg&country=6&maxPrice=0.5
```

### 取消激活
```
GET https://hero-sms.com/stubs/handler_api.php?action=cancelActivation&api_key=KEY&id=123456789
```

### 完成激活
```
GET https://hero-sms.com/stubs/handler_api.php?action=finishActivation&api_key=KEY&id=123456789
```
