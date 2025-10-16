# send requests to go service /enqueue endpoint
# acts as front-end/job submitter

from fastapi import FastAPI


app = FastAPI()


@app.get("/")
async def root():
    return {"message": "hello client app"}
