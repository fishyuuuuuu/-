const FINAL_FALLBACK_IMAGE = `data:image/svg+xml;utf8,${encodeURIComponent(
  '<svg xmlns="http://www.w3.org/2000/svg" width="400" height="400"><rect width="100%" height="100%" fill="#eef2ff"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="#6b7280" font-size="28" font-family="Arial">No Image</text></svg>'
)}`;

const IMAGE_BY_PRODUCT_ID = {
  1: 'https://images.unsplash.com/photo-1592286927505-1def25115558?w=400&h=400&fit=crop',
  2: 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=400&h=400&fit=crop',
  3: 'https://images.unsplash.com/photo-1610945265064-0e34e5519bbf?w=400&h=400&fit=crop',
  4: 'https://images.unsplash.com/photo-1580910051074-3eb694886505?w=400&h=400&fit=crop',
  5: 'https://images.unsplash.com/photo-1588423771073-b8903fbb85b5?w=400&h=400&fit=crop',
  6: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=400&h=400&fit=crop',
  7: 'https://images.unsplash.com/photo-1605901309584-818e25960a8f?w=400&h=400&fit=crop',
  8: 'https://images.unsplash.com/photo-1486401899868-0e435ed85128?w=400&h=400&fit=crop',
  9: 'https://images.unsplash.com/photo-1510127034890-ba27508e9f1c?w=400&h=400&fit=crop',
  10: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=400&h=400&fit=crop',
  11: 'https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=400&h=400&fit=crop',
  12: 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=400&h=400&fit=crop',
};

const POOL_BY_CATEGORY = {
  2: [
    'https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1610792516307-ea5acd9c3b00?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1619983081563-430f63602796?w=400&h=400&fit=crop'
  ],
  3: [
    'https://images.unsplash.com/photo-1515879218367-8466d910aaa4?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1593642702821-c8da6771f0c6?w=400&h=400&fit=crop'
  ],
  4: [
    'https://images.unsplash.com/photo-1556911220-bda9f7f7597e?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1586201375761-83865001e31c?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1581578731548-c64695cc6952?w=400&h=400&fit=crop'
  ],
  5: [
    'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1594633313593-bab3825d0caf?w=400&h=400&fit=crop'
  ],
  6: [
    'https://images.unsplash.com/photo-1571781926291-c477ebfd024b?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1556228720-195a672e8a03?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1598440947619-2c35fc9aa908?w=400&h=400&fit=crop'
  ],
  7: [
    'https://images.unsplash.com/photo-1517836357463-d25dfeac3438?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1571731956672-f2b94d7dd0cb?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1518611012118-696072aa579a?w=400&h=400&fit=crop'
  ],
  8: [
    'https://images.unsplash.com/photo-1546549032-9571cd6b27df?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1504674900247-0877df9cc836?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1498837167922-ddd27525d352?w=400&h=400&fit=crop'
  ],
  1: [
    'https://images.unsplash.com/photo-1483985988355-763728e1935b?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1542838132-92c53300491e?w=400&h=400&fit=crop',
    'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=400&h=400&fit=crop'
  ]
};

function isUsableRemoteImage(url) {
  return typeof url === 'string' && /^https?:\/\//.test(url);
}

function isPlaceholder(url) {
  if (typeof url !== 'string') return true;
  return (
    url.includes('via.placeholder.com') ||
    url.includes('placehold.co') ||
    url.startsWith('data:image/svg+xml')
  );
}

function hashString(input) {
  let hash = 0;
  for (let i = 0; i < input.length; i += 1) {
    hash = (hash << 5) - hash + input.charCodeAt(i);
    hash |= 0;
  }
  return Math.abs(hash);
}

function inferCategory(name = '', description = '') {
  const text = `${name} ${description}`.toLowerCase();
  if (/手机|iphone|huawei|xiaomi|samsung|耳机|相机|手表|switch|playstation/.test(text)) return 2;
  if (/电脑|笔记本|macbook|ipad|键盘|鼠标|显示器|显卡/.test(text)) return 3;
  if (/空调|冰箱|洗衣机|电视|吸尘器|扫地|家电/.test(text)) return 4;
  if (/鞋|服装|衣服|羽绒服|包/.test(text)) return 5;
  if (/护肤|美妆|口红|香水|精华|面霜/.test(text)) return 6;
  if (/运动|户外|瑜伽|健身/.test(text)) return 7;
  if (/食品|饮料|咖啡|零食|矿泉水|白酒/.test(text)) return 8;
  return 1;
}

function getStablePoolImage(product) {
  const category = Number(product?.category) || inferCategory(product?.name, product?.description);
  const pool = POOL_BY_CATEGORY[category] || POOL_BY_CATEGORY[1];
  const seedBase = Number(product?.id) || hashString(product?.name || 'product');
  return pool[seedBase % pool.length];
}

export function getProductImageCandidates(product) {
  const raw = product?.coverImage || product?.image || '';
  const normalizedRaw = typeof raw === 'string' && raw.startsWith('/')
    ? `${window.location.origin}${raw}`
    : raw;

  const candidates = [];
  if (isUsableRemoteImage(normalizedRaw) && !isPlaceholder(normalizedRaw)) {
    candidates.push(normalizedRaw);
  }

  const byId = IMAGE_BY_PRODUCT_ID[Number(product?.id)];
  if (byId) candidates.push(byId);

  candidates.push(getStablePoolImage(product));
  candidates.push(FINAL_FALLBACK_IMAGE);

  return [...new Set(candidates)];
}

export { FINAL_FALLBACK_IMAGE };
