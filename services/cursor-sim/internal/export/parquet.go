package export

// Parquet Export - DEFERRED
//
// Parquet export requires external dependencies (github.com/xitongsys/parquet-go)
// which add significant build complexity.
//
// For now, researchers can use:
// - CSV export (compatible with pandas, R, Excel)
// - JSON export (compatible with most data tools)
// - NDJSON streaming for large datasets
//
// If Parquet export is needed:
// 1. Add dependency: go get github.com/xitongsys/parquet-go
// 2. Implement ParquetExporter with WriteParquet method
// 3. Use parquet-go's writer interface
//
// Example implementation pattern:
//
//   type ParquetExporter struct {
//       writer parquet.Writer
//   }
//
//   func (e *ParquetExporter) WriteParquet(dataPoints []models.ResearchDataPoint) error {
//       for _, dp := range dataPoints {
//           if err := e.writer.Write(dp); err != nil {
//               return err
//           }
//       }
//       return e.writer.Close()
//   }
//
// For most research use cases, CSV or JSON will be sufficient as they can be
// easily converted to Parquet using Python/pandas:
//
//   import pandas as pd
//   df = pd.read_csv('research_data.csv')
//   df.to_parquet('research_data.parquet')
