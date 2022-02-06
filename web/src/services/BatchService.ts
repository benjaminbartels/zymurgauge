import { HTTP } from "../http-common";
import { BatchDetail, BatchSummary } from "../types/Batch";

class BatchService {
  getAllSummaries() {
    return HTTP.get<Array<BatchSummary>>("/batches");
  }
  getDetail(id: string) {
    return HTTP.get<BatchDetail>(`/batches/${id}`);
  }
}

export default new BatchService();
