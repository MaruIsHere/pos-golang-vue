export function openDB(name: string, version: number = 1): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(name, version);
    request.onupgradeneeded = (e) => {
      const db = request.result;
      if (!db.objectStoreNames.contains("products")) {
        db.createObjectStore("products", { keyPath: "id" });
      }
      if (!db.objectStoreNames.contains("offline_orders")) {
        db.createObjectStore("offline_orders", { keyPath: "id", autoIncrement: true });
      }
    };
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}

export async function setProducts(products: any[]) {
  const db = await openDB("pos_db");
  const tx = db.transaction("products", "readwrite");
  const store = tx.objectStore("products");
  store.clear();
  products.forEach(p => store.put(p));
  return new Promise((resolve, reject) => {
    tx.oncomplete = () => resolve(true);
    tx.onerror = () => reject(tx.error);
  });
}

export async function getProducts(): Promise<any[]> {
  const db = await openDB("pos_db");
  const tx = db.transaction("products", "readonly");
  const store = tx.objectStore("products");
  const request = store.getAll();
  return new Promise((resolve, reject) => {
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}

export async function addOfflineOrder(order: any): Promise<number> {
  const db = await openDB("pos_db");
  const tx = db.transaction("offline_orders", "readwrite");
  const store = tx.objectStore("offline_orders");
  const request = store.add({ ...order, timestamp: Date.now() });
  return new Promise((resolve, reject) => {
    request.onsuccess = () => resolve(request.result as number);
    request.onerror = () => reject(request.error);
  });
}

export async function getOfflineOrders(): Promise<any[]> {
  const db = await openDB("pos_db");
  const tx = db.transaction("offline_orders", "readonly");
  const store = tx.objectStore("offline_orders");
  const request = store.getAll();
  return new Promise((resolve, reject) => {
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}

export async function deleteOfflineOrder(id: number) {
  const db = await openDB("pos_db");
  const tx = db.transaction("offline_orders", "readwrite");
  const store = tx.objectStore("offline_orders");
  const request = store.delete(id);
  return new Promise((resolve, reject) => {
    request.onsuccess = () => resolve(true);
    request.onerror = () => reject(request.error);
  });
}
