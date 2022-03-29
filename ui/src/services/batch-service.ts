import axios from "axios";
import { BatchDetail, BatchSummary } from "../types/Batch";
import { authHeader, getUrl } from "./common";

class BatchService {
  getAllSummaries() {
    return axios.get<Array<BatchSummary>>(getUrl("batches"), {
      headers: authHeader(),
    });
  }
  getDetail(id: string) {
    return axios.get<BatchDetail>(getUrl(`batches/${id}`), {
      headers: authHeader(),
    });
  }
}

export default new BatchService();
