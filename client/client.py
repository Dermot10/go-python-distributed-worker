from typing import Optional, Dict, Any
from datetime import datetime
import asyncio
import random
import httpx


async def enqueue():
    async with httpx.AsyncClient() as client:

        for num in range(random.randint(1, 100)):
            created_at = datetime.utcnow().replace(microsecond=0).isoformat() + "Z"
            job = {
                "id": f"job-{num}",
                "type": "data_processing",
                "payload": {"filename": f"input-{num}.csv", "bucket": f"data-bucket-{num}"},
                "created_at": created_at
            }
            resp = await client.post("http://localhost:8081/enqueue", json=job)
            print(f"Sent {job['id']}, status: {resp.status_code}")
    return {"Status": "sent"}


if __name__ == "__main__":
    asyncio.run(enqueue())
