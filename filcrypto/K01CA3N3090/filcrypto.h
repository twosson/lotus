/* filcrypto Header */

#ifdef __cplusplus
extern "C" {
#endif


#ifndef filcrypto_H
#define filcrypto_H

/* Generated with cbindgen:0.16.0 */

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

#define SIGNATURE_BYTES 96

#define PRIVATE_KEY_BYTES 32

#define PUBLIC_KEY_BYTES 48

#define DIGEST_BYTES 96

typedef enum FCPResponseStatus {
  FCPResponseStatus_FCPNoError = 0,
  FCPResponseStatus_FCPUnclassifiedError = 1,
  FCPResponseStatus_FCPCallerError = 2,
  FCPResponseStatus_FCPReceiverError = 3,
} FCPResponseStatus;

typedef enum fil_RegisteredPoStProof {
  fil_RegisteredPoStProof_StackedDrgWinning2KiBV1,
  fil_RegisteredPoStProof_StackedDrgWinning8MiBV1,
  fil_RegisteredPoStProof_StackedDrgWinning512MiBV1,
  fil_RegisteredPoStProof_StackedDrgWinning32GiBV1,
  fil_RegisteredPoStProof_StackedDrgWinning64GiBV1,
  fil_RegisteredPoStProof_StackedDrgWindow2KiBV1,
  fil_RegisteredPoStProof_StackedDrgWindow8MiBV1,
  fil_RegisteredPoStProof_StackedDrgWindow512MiBV1,
  fil_RegisteredPoStProof_StackedDrgWindow32GiBV1,
  fil_RegisteredPoStProof_StackedDrgWindow64GiBV1,
} fil_RegisteredPoStProof;

typedef enum fil_RegisteredSealProof {
  fil_RegisteredSealProof_StackedDrg2KiBV1,
  fil_RegisteredSealProof_StackedDrg8MiBV1,
  fil_RegisteredSealProof_StackedDrg512MiBV1,
  fil_RegisteredSealProof_StackedDrg32GiBV1,
  fil_RegisteredSealProof_StackedDrg64GiBV1,
  fil_RegisteredSealProof_StackedDrg2KiBV1_1,
  fil_RegisteredSealProof_StackedDrg8MiBV1_1,
  fil_RegisteredSealProof_StackedDrg512MiBV1_1,
  fil_RegisteredSealProof_StackedDrg32GiBV1_1,
  fil_RegisteredSealProof_StackedDrg64GiBV1_1,
} fil_RegisteredSealProof;

typedef struct fil_BLSDigest {
  uint8_t inner[DIGEST_BYTES];
} fil_BLSDigest;

/**
 * HashResponse
 */
typedef struct fil_HashResponse {
  struct fil_BLSDigest digest;
} fil_HashResponse;

typedef struct fil_BLSSignature {
  uint8_t inner[SIGNATURE_BYTES];
} fil_BLSSignature;

/**
 * AggregateResponse
 */
typedef struct fil_AggregateResponse {
  struct fil_BLSSignature signature;
} fil_AggregateResponse;

typedef struct fil_BLSPrivateKey {
  uint8_t inner[PRIVATE_KEY_BYTES];
} fil_BLSPrivateKey;

/**
 * PrivateKeyGenerateResponse
 */
typedef struct fil_PrivateKeyGenerateResponse {
  struct fil_BLSPrivateKey private_key;
} fil_PrivateKeyGenerateResponse;

typedef struct fil_32ByteArray {
  uint8_t inner[32];
} fil_32ByteArray;

/**
 * PrivateKeySignResponse
 */
typedef struct fil_PrivateKeySignResponse {
  struct fil_BLSSignature signature;
} fil_PrivateKeySignResponse;

typedef struct fil_BLSPublicKey {
  uint8_t inner[PUBLIC_KEY_BYTES];
} fil_BLSPublicKey;

/**
 * PrivateKeyPublicKeyResponse
 */
typedef struct fil_PrivateKeyPublicKeyResponse {
  struct fil_BLSPublicKey public_key;
} fil_PrivateKeyPublicKeyResponse;

typedef struct fil_WriteWithAlignmentResponse {
  uint8_t comm_p[32];
  const char *error_msg;
  uint64_t left_alignment_unpadded;
  enum FCPResponseStatus status_code;
  uint64_t total_write_unpadded;
} fil_WriteWithAlignmentResponse;

typedef struct fil_WriteWithoutAlignmentResponse {
  uint8_t comm_p[32];
  const char *error_msg;
  enum FCPResponseStatus status_code;
  uint64_t total_write_unpadded;
} fil_WriteWithoutAlignmentResponse;

typedef struct fil_FauxRepResponse {
  const char *error_msg;
  enum FCPResponseStatus status_code;
  uint8_t commitment[32];
} fil_FauxRepResponse;

typedef struct fil_SealPreCommitPhase1Response {
  const char *error_msg;
  enum FCPResponseStatus status_code;
  const uint8_t *seal_pre_commit_phase1_output_ptr;
  size_t seal_pre_commit_phase1_output_len;
} fil_SealPreCommitPhase1Response;

typedef struct fil_PublicPieceInfo {
  uint64_t num_bytes;
  uint8_t comm_p[32];
} fil_PublicPieceInfo;

typedef struct fil_SealPreCommitPhase2Response {
  const char *error_msg;
  enum FCPResponseStatus status_code;
  enum fil_RegisteredSealProof registered_proof;
  uint8_t comm_d[32];
  uint8_t comm_r[32];
} fil_SealPreCommitPhase2Response;

typedef struct fil_SealCommitPhase1Response {
  enum FCPResponseStatus status_code;
  const char *error_msg;
  const uint8_t *seal_commit_phase1_output_ptr;
  size_t seal_commit_phase1_output_len;
} fil_SealCommitPhase1Response;

typedef struct fil_SealCommitPhase2Response {
  enum FCPResponseStatus status_code;
  const char *error_msg;
  const uint8_t *proof_ptr;
  size_t proof_len;
} fil_SealCommitPhase2Response;

typedef struct fil_UnsealRangeResponse {
  enum FCPResponseStatus status_code;
  const char *error_msg;
} fil_UnsealRangeResponse;

typedef struct fil_VerifySealResponse {
  enum FCPResponseStatus status_code;
  const char *error_msg;
  bool is_valid;
} fil_VerifySealResponse;

typedef struct fil_GenerateWinningPoStSectorChallenge {
  const char *error_msg;
  enum FCPResponseStatus status_code;
  const uint64_t *ids_ptr;
  size_t ids_len;
} fil_GenerateWinningPoStSectorChallenge;

typedef struct fil_GenerateFallbackSectorChallengesResponse {
  const char *error_msg;
  enum FCPResponseStatus status_code;
  const uint64_t *ids_ptr;
  size_t ids_len;
  const uint64_t *challenges_ptr;
  size_t challenges_len;
  size_t challenges_stride;
} fil_GenerateFallbackSectorChallengesResponse;

typedef struct fil_VanillaProof {
  size_t proof_len;
  const uint8_t *proof_ptr;
} fil_VanillaProof;

typedef struct fil_GenerateSingleVanillaProofResponse {
  const char *error_msg;
  struct fil_VanillaProof vanilla_proof;
  enum FCPResponseStatus status_code;
} fil_GenerateSingleVanillaProofResponse;

typedef struct fil_PrivateReplicaInfo {
  enum fil_RegisteredPoStProof registered_proof;
  const char *cache_dir_path;
  uint8_t comm_r[32];
  const char *replica_path;
  uint64_t sector_id;
} fil_PrivateReplicaInfo;

typedef struct fil_PoStProof {
  enum fil_RegisteredPoStProof registered_proof;
  size_t proof_len;
  const uint8_t *proof_ptr;
} fil_PoStProof;

typedef struct fil_GenerateWinningPoStResponse {
  const char *error_msg;
  size_t proofs_len;
  const struct fil_PoStProof *proofs_ptr;
  enum FCPResponseStatus status_code;
} fil_GenerateWinningPoStResponse;

typedef struct fil_VerifyWinningPoStResponse {
  enum FCPResponseStatus status_code;
  const char *error_msg;
  bool is_valid;
} fil_VerifyWinningPoStResponse;

typedef struct fil_PublicReplicaInfo {
  enum fil_RegisteredPoStProof registered_proof;
  uint8_t comm_r[32];
  uint64_t sector_id;
} fil_PublicReplicaInfo;

typedef struct fil_GenerateWindowPoStResponse {
  const char *error_msg;
  size_t proofs_len;
  const struct fil_PoStProof *proofs_ptr;
  size_t faulty_sectors_len;
  const uint64_t *faulty_sectors_ptr;
  enum FCPResponseStatus status_code;
} fil_GenerateWindowPoStResponse;

typedef struct fil_VerifyWindowPoStResponse {
  enum FCPResponseStatus status_code;
  const char *error_msg;
  bool is_valid;
} fil_VerifyWindowPoStResponse;

typedef struct fil_GeneratePieceCommitmentResponse {
  enum FCPResponseStatus status_code;
  const char *error_msg;
  uint8_t comm_p[32];
  /**
   * The number of unpadded bytes in the original piece plus any (unpadded)
   * alignment bytes added to create a whole merkle tree.
   */
  uint64_t num_bytes_aligned;
} fil_GeneratePieceCommitmentResponse;

typedef struct fil_GenerateDataCommitmentResponse {
  enum FCPResponseStatus status_code;
  const char *error_msg;
  uint8_t comm_d[32];
} fil_GenerateDataCommitmentResponse;

typedef struct fil_ClearCacheResponse {
  const char *error_msg;
  enum FCPResponseStatus status_code;
} fil_ClearCacheResponse;

/**
 *
 */
typedef struct fil_StringResponse {
  enum FCPResponseStatus status_code;
  const char *error_msg;
  const char *string_val;
} fil_StringResponse;

typedef struct fil_FinalizeTicketResponse {
  enum FCPResponseStatus status_code;
  const char *error_msg;
  uint8_t ticket[32];
} fil_FinalizeTicketResponse;

typedef struct fil_GpuDeviceResponse {
  enum FCPResponseStatus status_code;
  const char *error_msg;
  size_t devices_len;
  const char *const *devices_ptr;
} fil_GpuDeviceResponse;

typedef struct fil_InitLogFdResponse {
  enum FCPResponseStatus status_code;
  const char *error_msg;
} fil_InitLogFdResponse;

/**
 * Compute the digest of a message
 *
 * # Arguments
 *
 * * `message_ptr` - pointer to a message byte array
 * * `message_len` - length of the byte array
 */
struct fil_HashResponse *fil_hash(const uint8_t *message_ptr, size_t message_len);

/**
 * Aggregate signatures together into a new signature
 *
 * # Arguments
 *
 * * `flattened_signatures_ptr` - pointer to a byte array containing signatures
 * * `flattened_signatures_len` - length of the byte array (multiple of SIGNATURE_BYTES)
 *
 * Returns `NULL` on error. Result must be freed using `destroy_aggregate_response`.
 */
struct fil_AggregateResponse *fil_aggregate(const uint8_t *flattened_signatures_ptr,
                                            size_t flattened_signatures_len);

/**
 * Verify that a signature is the aggregated signature of hashes - pubkeys
 *
 * # Arguments
 *
 * * `signature_ptr`             - pointer to a signature byte array (SIGNATURE_BYTES long)
 * * `flattened_digests_ptr`     - pointer to a byte array containing digests
 * * `flattened_digests_len`     - length of the byte array (multiple of DIGEST_BYTES)
 * * `flattened_public_keys_ptr` - pointer to a byte array containing public keys
 * * `flattened_public_keys_len` - length of the array
 */
int fil_verify(const uint8_t *signature_ptr,
               const uint8_t *flattened_digests_ptr,
               size_t flattened_digests_len,
               const uint8_t *flattened_public_keys_ptr,
               size_t flattened_public_keys_len);

/**
 * Verify that a signature is the aggregated signature of the hhashed messages
 *
 * # Arguments
 *
 * * `signature_ptr`             - pointer to a signature byte array (SIGNATURE_BYTES long)
 * * `messages_ptr`              - pointer to an array containing the pointers to the messages
 * * `messages_sizes_ptr`        - pointer to an array containing the lengths of the messages
 * * `messages_len`              - length of the two messages arrays
 * * `flattened_public_keys_ptr` - pointer to a byte array containing public keys
 * * `flattened_public_keys_len` - length of the array
 */
int fil_hash_verify(const uint8_t *signature_ptr,
                    const uint8_t *flattened_messages_ptr,
                    size_t flattened_messages_len,
                    const size_t *message_sizes_ptr,
                    size_t message_sizes_len,
                    const uint8_t *flattened_public_keys_ptr,
                    size_t flattened_public_keys_len);

/**
 * Generate a new private key
 */
struct fil_PrivateKeyGenerateResponse *fil_private_key_generate(void);

/**
 * Generate a new private key with seed
 *
 * **Warning**: Use this function only for testing or with very secure seeds
 *
 * # Arguments
 *
 * * `raw_seed` - a seed byte array with 32 bytes
 *
 * Returns `NULL` when passed a NULL pointer.
 */
struct fil_PrivateKeyGenerateResponse *fil_private_key_generate_with_seed(struct fil_32ByteArray raw_seed);

/**
 * Sign a message with a private key and return the signature
 *
 * # Arguments
 *
 * * `raw_private_key_ptr` - pointer to a private key byte array
 * * `message_ptr` - pointer to a message byte array
 * * `message_len` - length of the byte array
 *
 * Returns `NULL` when passed invalid arguments.
 */
struct fil_PrivateKeySignResponse *fil_private_key_sign(const uint8_t *raw_private_key_ptr,
                                                        const uint8_t *message_ptr,
                                                        size_t message_len);

/**
 * Generate the public key for a private key
 *
 * # Arguments
 *
 * * `raw_private_key_ptr` - pointer to a private key byte array
 *
 * Returns `NULL` when passed invalid arguments.
 */
struct fil_PrivateKeyPublicKeyResponse *fil_private_key_public_key(const uint8_t *raw_private_key_ptr);

void fil_destroy_hash_response(struct fil_HashResponse *ptr);

void fil_destroy_aggregate_response(struct fil_AggregateResponse *ptr);

void fil_destroy_private_key_generate_response(struct fil_PrivateKeyGenerateResponse *ptr);

void fil_destroy_private_key_sign_response(struct fil_PrivateKeySignResponse *ptr);

void fil_destroy_private_key_public_key_response(struct fil_PrivateKeyPublicKeyResponse *ptr);

/**
 * TODO: document
 *
 */
struct fil_WriteWithAlignmentResponse *fil_write_with_alignment(enum fil_RegisteredSealProof registered_proof,
                                                                int src_fd,
                                                                uint64_t src_size,
                                                                int dst_fd,
                                                                const uint64_t *existing_piece_sizes_ptr,
                                                                size_t existing_piece_sizes_len);

/**
 * TODO: document
 *
 */
struct fil_WriteWithoutAlignmentResponse *fil_write_without_alignment(enum fil_RegisteredSealProof registered_proof,
                                                                      int src_fd,
                                                                      uint64_t src_size,
                                                                      int dst_fd);

struct fil_FauxRepResponse *fil_fauxrep(enum fil_RegisteredSealProof registered_proof,
                                        const char *cache_dir_path,
                                        const char *sealed_sector_path);

struct fil_FauxRepResponse *fil_fauxrep2(enum fil_RegisteredSealProof registered_proof,
                                         const char *cache_dir_path,
                                         const char *existing_p_aux_path);

/**
 * TODO: document
 *
 */
struct fil_SealPreCommitPhase1Response *fil_seal_pre_commit_phase1(enum fil_RegisteredSealProof registered_proof,
                                                                   const char *cache_dir_path,
                                                                   const char *staged_sector_path,
                                                                   const char *sealed_sector_path,
                                                                   uint64_t sector_id,
                                                                   struct fil_32ByteArray prover_id,
                                                                   struct fil_32ByteArray ticket,
                                                                   const struct fil_PublicPieceInfo *pieces_ptr,
                                                                   size_t pieces_len);

/**
 * TODO: document
 *
 */
struct fil_SealPreCommitPhase2Response *fil_seal_pre_commit_phase2(const uint8_t *seal_pre_commit_phase1_output_ptr,
                                                                   size_t seal_pre_commit_phase1_output_len,
                                                                   const char *cache_dir_path,
                                                                   const char *sealed_sector_path);

/**
 * TODO: document
 *
 */
struct fil_SealCommitPhase1Response *fil_seal_commit_phase1(enum fil_RegisteredSealProof registered_proof,
                                                            struct fil_32ByteArray comm_r,
                                                            struct fil_32ByteArray comm_d,
                                                            const char *cache_dir_path,
                                                            const char *replica_path,
                                                            uint64_t sector_id,
                                                            struct fil_32ByteArray prover_id,
                                                            struct fil_32ByteArray ticket,
                                                            struct fil_32ByteArray seed,
                                                            const struct fil_PublicPieceInfo *pieces_ptr,
                                                            size_t pieces_len);

struct fil_SealCommitPhase2Response *fil_seal_commit_phase2(const uint8_t *seal_commit_phase1_output_ptr,
                                                            size_t seal_commit_phase1_output_len,
                                                            uint64_t sector_id,
                                                            struct fil_32ByteArray prover_id);

/**
 * TODO: document
 */
struct fil_UnsealRangeResponse *fil_unseal_range(enum fil_RegisteredSealProof registered_proof,
                                                 const char *cache_dir_path,
                                                 int sealed_sector_fd_raw,
                                                 int unseal_output_fd_raw,
                                                 uint64_t sector_id,
                                                 struct fil_32ByteArray prover_id,
                                                 struct fil_32ByteArray ticket,
                                                 struct fil_32ByteArray comm_d,
                                                 uint64_t unpadded_byte_index,
                                                 uint64_t unpadded_bytes_amount);

/**
 * Verifies the output of seal.
 *
 */
struct fil_VerifySealResponse *fil_verify_seal(enum fil_RegisteredSealProof registered_proof,
                                               struct fil_32ByteArray comm_r,
                                               struct fil_32ByteArray comm_d,
                                               struct fil_32ByteArray prover_id,
                                               struct fil_32ByteArray ticket,
                                               struct fil_32ByteArray seed,
                                               uint64_t sector_id,
                                               const uint8_t *proof_ptr,
                                               size_t proof_len);

/**
 * TODO: document
 *
 */
struct fil_GenerateWinningPoStSectorChallenge *fil_generate_winning_post_sector_challenge(enum fil_RegisteredPoStProof registered_proof,
                                                                                          struct fil_32ByteArray randomness,
                                                                                          uint64_t sector_set_len,
                                                                                          struct fil_32ByteArray prover_id);

/**
 * TODO: document
 *
 */
struct fil_GenerateFallbackSectorChallengesResponse *fil_generate_fallback_sector_challenges(enum fil_RegisteredPoStProof registered_proof,
                                                                                             struct fil_32ByteArray randomness,
                                                                                             const uint64_t *sector_ids_ptr,
                                                                                             size_t sector_ids_len,
                                                                                             struct fil_32ByteArray prover_id);

/**
 * TODO: document
 *
 */
struct fil_GenerateSingleVanillaProofResponse *fil_generate_single_vanilla_proof(struct fil_PrivateReplicaInfo replica,
                                                                                 const uint64_t *challenges_ptr,
                                                                                 size_t challenges_len);

/**
 * TODO: document
 *
 */
struct fil_GenerateWinningPoStResponse *fil_generate_winning_post_with_vanilla(enum fil_RegisteredPoStProof registered_proof,
                                                                               struct fil_32ByteArray randomness,
                                                                               struct fil_32ByteArray prover_id,
                                                                               const struct fil_VanillaProof *vanilla_proofs_ptr,
                                                                               size_t vanilla_proofs_len);

/**
 * TODO: document
 *
 */
struct fil_GenerateWinningPoStResponse *fil_generate_winning_post(struct fil_32ByteArray randomness,
                                                                  const struct fil_PrivateReplicaInfo *replicas_ptr,
                                                                  size_t replicas_len,
                                                                  struct fil_32ByteArray prover_id);

/**
 * Verifies that a proof-of-spacetime is valid.
 */
struct fil_VerifyWinningPoStResponse *fil_verify_winning_post(struct fil_32ByteArray randomness,
                                                              const struct fil_PublicReplicaInfo *replicas_ptr,
                                                              size_t replicas_len,
                                                              const struct fil_PoStProof *proofs_ptr,
                                                              size_t proofs_len,
                                                              struct fil_32ByteArray prover_id);

/**
 * TODO: document
 *
 */
struct fil_GenerateWindowPoStResponse *fil_generate_window_post_with_vanilla(enum fil_RegisteredPoStProof registered_proof,
                                                                             struct fil_32ByteArray randomness,
                                                                             struct fil_32ByteArray prover_id,
                                                                             const struct fil_VanillaProof *vanilla_proofs_ptr,
                                                                             size_t vanilla_proofs_len);

/**
 * TODO: document
 *
 */
struct fil_GenerateWindowPoStResponse *fil_generate_window_post(struct fil_32ByteArray randomness,
                                                                const struct fil_PrivateReplicaInfo *replicas_ptr,
                                                                size_t replicas_len,
                                                                struct fil_32ByteArray prover_id);

/**
 * Verifies that a proof-of-spacetime is valid.
 */
struct fil_VerifyWindowPoStResponse *fil_verify_window_post(struct fil_32ByteArray randomness,
                                                            const struct fil_PublicReplicaInfo *replicas_ptr,
                                                            size_t replicas_len,
                                                            const struct fil_PoStProof *proofs_ptr,
                                                            size_t proofs_len,
                                                            struct fil_32ByteArray prover_id);

/**
 * Returns the merkle root for a piece after piece padding and alignment.
 * The caller is responsible for closing the passed in file descriptor.
 */
struct fil_GeneratePieceCommitmentResponse *fil_generate_piece_commitment(enum fil_RegisteredSealProof registered_proof,
                                                                          int piece_fd_raw,
                                                                          uint64_t unpadded_piece_size);

/**
 * Returns the merkle root for a sector containing the provided pieces.
 */
struct fil_GenerateDataCommitmentResponse *fil_generate_data_commitment(enum fil_RegisteredSealProof registered_proof,
                                                                        const struct fil_PublicPieceInfo *pieces_ptr,
                                                                        size_t pieces_len);

struct fil_ClearCacheResponse *fil_clear_cache(uint64_t sector_size, const char *cache_dir_path);

void fil_destroy_write_with_alignment_response(struct fil_WriteWithAlignmentResponse *ptr);

void fil_destroy_write_without_alignment_response(struct fil_WriteWithoutAlignmentResponse *ptr);

void fil_destroy_fauxrep_response(struct fil_FauxRepResponse *ptr);

void fil_destroy_seal_pre_commit_phase1_response(struct fil_SealPreCommitPhase1Response *ptr);

void fil_destroy_seal_pre_commit_phase2_response(struct fil_SealPreCommitPhase2Response *ptr);

void fil_destroy_seal_commit_phase1_response(struct fil_SealCommitPhase1Response *ptr);

void fil_destroy_seal_commit_phase2_response(struct fil_SealCommitPhase2Response *ptr);

void fil_destroy_unseal_range_response(struct fil_UnsealRangeResponse *ptr);

void fil_destroy_generate_piece_commitment_response(struct fil_GeneratePieceCommitmentResponse *ptr);

void fil_destroy_generate_data_commitment_response(struct fil_GenerateDataCommitmentResponse *ptr);

void fil_destroy_string_response(struct fil_StringResponse *ptr);

/**
 * Returns the number of user bytes that will fit into a staged sector.
 *
 */
uint64_t fil_get_max_user_bytes_per_staged_sector(enum fil_RegisteredSealProof registered_proof);

/**
 * Returns the CID of the Groth parameter file for sealing.
 *
 */
struct fil_StringResponse *fil_get_seal_params_cid(enum fil_RegisteredSealProof registered_proof);

/**
 * Returns the CID of the verifying key-file for verifying a seal proof.
 *
 */
struct fil_StringResponse *fil_get_seal_verifying_key_cid(enum fil_RegisteredSealProof registered_proof);

/**
 * Returns the path from which the proofs library expects to find the Groth
 * parameter file used when sealing.
 *
 */
struct fil_StringResponse *fil_get_seal_params_path(enum fil_RegisteredSealProof registered_proof);

/**
 * Returns the path from which the proofs library expects to find the verifying
 * key-file used when verifying a seal proof.
 *
 */
struct fil_StringResponse *fil_get_seal_verifying_key_path(enum fil_RegisteredSealProof registered_proof);

/**
 * Returns the identity of the circuit for the provided seal proof.
 *
 */
struct fil_StringResponse *fil_get_seal_circuit_identifier(enum fil_RegisteredSealProof registered_proof);

/**
 * Returns the version of the provided seal proof type.
 *
 */
struct fil_StringResponse *fil_get_seal_version(enum fil_RegisteredSealProof registered_proof);

/**
 * Returns the CID of the Groth parameter file for generating a PoSt.
 *
 */
struct fil_StringResponse *fil_get_post_params_cid(enum fil_RegisteredPoStProof registered_proof);

/**
 * Returns the CID of the verifying key-file for verifying a PoSt proof.
 *
 */
struct fil_StringResponse *fil_get_post_verifying_key_cid(enum fil_RegisteredPoStProof registered_proof);

/**
 * Returns the path from which the proofs library expects to find the Groth
 * parameter file used when generating a PoSt.
 *
 */
struct fil_StringResponse *fil_get_post_params_path(enum fil_RegisteredPoStProof registered_proof);

/**
 * Returns the path from which the proofs library expects to find the verifying
 * key-file used when verifying a PoSt proof.
 *
 */
struct fil_StringResponse *fil_get_post_verifying_key_path(enum fil_RegisteredPoStProof registered_proof);

/**
 * Returns the identity of the circuit for the provided PoSt proof type.
 *
 */
struct fil_StringResponse *fil_get_post_circuit_identifier(enum fil_RegisteredPoStProof registered_proof);

/**
 * Returns the version of the provided seal proof.
 *
 */
struct fil_StringResponse *fil_get_post_version(enum fil_RegisteredPoStProof registered_proof);

/**
 * Deallocates a VerifySealResponse.
 *
 */
void fil_destroy_verify_seal_response(struct fil_VerifySealResponse *ptr);

void fil_destroy_finalize_ticket_response(struct fil_FinalizeTicketResponse *ptr);

/**
 * Deallocates a VerifyPoStResponse.
 *
 */
void fil_destroy_verify_winning_post_response(struct fil_VerifyWinningPoStResponse *ptr);

void fil_destroy_verify_window_post_response(struct fil_VerifyWindowPoStResponse *ptr);

void fil_destroy_generate_fallback_sector_challenges_response(struct fil_GenerateFallbackSectorChallengesResponse *ptr);

void fil_destroy_generate_single_vanilla_proof_response(struct fil_GenerateSingleVanillaProofResponse *ptr);

void fil_destroy_generate_winning_post_response(struct fil_GenerateWinningPoStResponse *ptr);

void fil_destroy_generate_window_post_response(struct fil_GenerateWindowPoStResponse *ptr);

void fil_destroy_generate_winning_post_sector_challenge(struct fil_GenerateWinningPoStSectorChallenge *ptr);

void fil_destroy_clear_cache_response(struct fil_ClearCacheResponse *ptr);

/**
 * Returns an array of strings containing the device names that can be used.
 */
struct fil_GpuDeviceResponse *fil_get_gpu_devices(void);

/**
 * Initializes the logger with a file descriptor where logs will be logged into.
 *
 * This is usually a pipe that was opened on the receiving side of the logs. The logger is
 * initialized on the invocation, subsequent calls won't have any effect.
 *
 * This function must be called right at the start, before any other call. Else the logger will
 * be initializes implicitely and log to stderr.
 */
struct fil_InitLogFdResponse *fil_init_log_fd(int log_fd);

void fil_destroy_gpu_device_response(struct fil_GpuDeviceResponse *ptr);

void fil_destroy_init_log_fd_response(struct fil_InitLogFdResponse *ptr);

#endif /* filcrypto_H */

#ifdef __cplusplus
} /* extern "C" */
#endif
