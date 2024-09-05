import axios from "axios";
import { useEffect, useState } from "react";

const useGetData = <T>(uri: string, mapper?: (obj: any) => T) => {
  const [data, setData] = useState<T | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<any>(null);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const response = await axios.get("http://localhost:4000" + uri, {
        withCredentials: true, // This will include cookies in the request
      });

      // Apply mapper if provided
      const mappedData = mapper
        ? mapper(response.data?.Datas)
        : response.data?.Datas;
      setData(mappedData);

      if (response.data?.Errors) {
        setError(response.data?.Errors);
      }
    } catch (err) {
      setError(err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [uri]);

  return {
    expect: data,
    error,
    errorPerso: error,
    isLoading,
    refetch: fetchData,
  };
};

export default useGetData;
